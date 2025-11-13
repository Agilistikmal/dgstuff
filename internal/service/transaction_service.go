package service

import (
	"context"
	"strconv"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment/xendit_payment"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TransactionService struct {
	db            *gorm.DB
	validator     *app.Validator
	xenditPayment *xendit_payment.XenditPayment
}

func NewTransactionService(db *gorm.DB, validator *app.Validator, xenditPayment *xendit_payment.XenditPayment) *TransactionService {
	return &TransactionService{db: db, validator: validator, xenditPayment: xenditPayment}
}

func (s *TransactionService) Create(ctx context.Context, dto model.TransactionCreateDTO) (*model.Transaction, error) {
	err := s.validator.Validate(dto)
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	id := pkg.GenerateRandomID(viper.GetString("transaction.prefix"), 10)
	transaction := model.Transaction{
		ID:       id,
		Email:    dto.Email,
		Currency: model.Currency(dto.Currency),
		Status:   model.TransactionStatus(dto.Status),
	}

	for _, stuff := range dto.Stuffs {
		var findStuff model.Stuff
		err = tx.Where("id = ?", stuff.StuffID).First(&findStuff).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("gorm failed to find stuff by id: %v", err)
			return nil, app.NewInternalServerError()
		}
		transaction.Stuffs = append(transaction.Stuffs, model.TransactionStuff{
			StuffID:    stuff.StuffID,
			Quantity:   stuff.Quantity,
			StuffName:  findStuff.Name,
			StuffSlug:  findStuff.Slug,
			StuffPrice: findStuff.Price,
			Currency:   findStuff.Currency,
			TotalPrice: findStuff.Price * float64(stuff.Quantity),
			Data:       nil,
		})
		transaction.Amount += findStuff.Price * float64(stuff.Quantity)
	}

	payment, err := s.createPayment(ctx, transaction, dto.PaymentProvider)
	if err != nil {
		tx.Rollback()
		logrus.Errorf("failed to create payment: %v", err)
		return nil, app.NewInternalServerError()
	}
	transaction.Payment = model.TransactionPayment{
		ID:            payment.ID,
		TransactionID: transaction.ID,
		Type:          payment.Type,
		Method:        payment.Method,
		Code:          payment.Code,
		Status:        payment.Status,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Provider:      payment.Provider,
		URL:           payment.URL,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}

	err = tx.Save(&transaction).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to create transaction: %v", err)
		return nil, app.NewInternalServerError()
	}

	return &transaction, nil
}

func (s *TransactionService) createPayment(ctx context.Context, transaction model.Transaction, paymentProvider payment.PaymentProvider) (*payment.Payment, error) {
	invoiceRequest := payment.PaymentInvoiceRequest{
		TransactionID: transaction.ID,
		Amount:        transaction.Amount,
		Currency:      string(transaction.Currency),
		Customer: payment.PaymentCustomer{
			Name:  transaction.Email,
			Email: transaction.Email,
		},
	}

	for _, stuff := range transaction.Stuffs {
		invoiceRequest.Items = append(invoiceRequest.Items, payment.PaymentInvoiceItem{
			ID:       strconv.Itoa(stuff.ID),
			Name:     stuff.StuffName,
			Quantity: stuff.Quantity,
			Price:    stuff.StuffPrice,
		})
	}

	switch paymentProvider {
	case payment.PaymentProviderXendit:
		return s.xenditPayment.CreateInvoice(ctx, invoiceRequest)
	default:
		return nil, app.NewInternalServerError()
	}
}
