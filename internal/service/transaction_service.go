package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TransactionService struct {
	db            *gorm.DB
	validator     *app.Validator
	xenditPayment payment.PaymentService
}

func NewTransactionService(db *gorm.DB, validator *app.Validator, xenditPayment payment.PaymentService) *TransactionService {
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
		Status:   model.TransactionStatusPending,
	}

	for _, stuff := range dto.Stuffs {
		var findStuff model.Stuff
		err = tx.Where("id = ?", stuff.StuffID).First(&findStuff).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("gorm failed to find stuff by id: %v", err)
			return nil, app.NewInternalServerError()
		}

		var findStock model.Stock
		err = tx.Where("stuff_id = ?", stuff.StuffID).First(&findStock).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("gorm failed to find stock by stuff id: %v", err)
			return nil, app.NewInternalServerError()
		}

		decryptedStockValues, err := pkg.Decrypt(findStock.Values, viper.GetString("stock.secret_key"))
		if err != nil {
			tx.Rollback()
			logrus.Errorf("failed to decrypt stock values: %v", err)
			return nil, app.NewInternalServerError()
		}
		stockValues := strings.Split(decryptedStockValues, findStock.Separator)
		if len(stockValues) < stuff.Quantity {
			tx.Rollback()
			logrus.Errorf("stock not enough: %v", stuff.Quantity)
			return nil, app.NewBadRequestError("stock not enough")
		}

		selectedStockValue := stockValues[0:stuff.Quantity]
		encryptedSelectedStockValue, err := pkg.Encrypt(strings.Join(selectedStockValue, findStock.Separator), viper.GetString("stock.secret_key"))
		if err != nil {
			tx.Rollback()
			logrus.Errorf("failed to encrypt selected stock value: %v", err)
			return nil, app.NewInternalServerError()
		}
		transactionStuffData := model.TransactionStuffData{
			Values:    encryptedSelectedStockValue,
			Separator: findStock.Separator,
		}
		transaction.Stuffs = append(transaction.Stuffs, model.TransactionStuff{
			StuffID:    stuff.StuffID,
			Quantity:   stuff.Quantity,
			StuffName:  findStuff.Name,
			StuffSlug:  findStuff.Slug,
			StuffPrice: findStuff.Price,
			Currency:   findStuff.Currency,
			TotalPrice: findStuff.Price * float64(stuff.Quantity),
			Data:       &transactionStuffData,
		})
		transaction.Amount += findStuff.Price * float64(stuff.Quantity)

		remainingStockValues := stockValues[stuff.Quantity:]
		encryptedRemainingStockValues, err := pkg.Encrypt(strings.Join(remainingStockValues, findStock.Separator), viper.GetString("stock.secret_key"))
		if err != nil {
			tx.Rollback()
			logrus.Errorf("failed to encrypt remaining stock values: %v", err)
			return nil, app.NewInternalServerError()
		}
		findStock.Count -= stuff.Quantity
		findStock.Values = encryptedRemainingStockValues
		err = tx.Save(&findStock).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorf("gorm failed to update stock: %v", err)
			return nil, app.NewInternalServerError()
		}
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

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to commit transaction: %v", err)
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
			Name:    transaction.Email,
			Email:   transaction.Email,
			Phone:   "081234567890",
			Address: "Jl. Raya No. 123",
			City:    "Jakarta",
			State:   "DKI Jakarta",
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

func (s *TransactionService) Get(ctx context.Context, transactionID string, token string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := s.db.Preload("Payment").Preload("Stuffs.Data").Where("id = ?", transactionID).First(&transaction).Error
	if err != nil {
		return nil, err
	}

	if transaction.Payment.Status != payment.PaymentStatusSuccess {
		var p *payment.Payment

		switch transaction.Payment.Provider {
		case payment.PaymentProviderXendit:
			p, err = s.xenditPayment.GetPayment(ctx, transaction.Payment.ID)
			if err != nil {
				return nil, err
			}
		default:
			logrus.Errorf("payment provider not supported: %s", transaction.Payment.Provider)
			return nil, app.NewBadRequestError("payment provider not supported")
		}

		if p == nil {
			logrus.Errorf("payment not found: %s", transaction.Payment.ID)
			return nil, app.NewNotFoundError("payment not found")
		}

		if p.Status != transaction.Payment.Status {
			transaction.Payment.Status = p.Status
			transaction.Payment.UpdatedAt = p.UpdatedAt

			switch p.Status {
			case payment.PaymentStatusSuccess:
				transaction.Status = model.TransactionStatusSuccess
			case payment.PaymentStatusFailed:
				transaction.Status = model.TransactionStatusFailed
			case payment.PaymentStatusPending:
				transaction.Status = model.TransactionStatusPending
			default:
				logrus.Errorf("payment status not supported: %s", p.Status)
				return nil, app.NewBadRequestError("payment status not supported")
			}
		}

		transaction.Payment = model.TransactionPayment{
			ID:            p.ID,
			TransactionID: p.TransactionID,
			Type:          p.Type,
			Method:        p.Method,
			Code:          p.Code,
			Status:        p.Status,
			Amount:        p.Amount,
			Currency:      p.Currency,
			Provider:      p.Provider,
			URL:           p.URL,
			ExpiresAt:     p.ExpiresAt,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}

		err = s.db.Save(&transaction).Error
		if err != nil {
			logrus.Errorf("gorm failed to update transaction: %v", err)
			return nil, app.NewInternalServerError()
		}
	}

	if transaction.Payment.Status != payment.PaymentStatusSuccess {
		var newStuffs []model.TransactionStuff
		for _, stuff := range transaction.Stuffs {
			stuff.Data = nil
			newStuffs = append(newStuffs, stuff)
		}
		transaction.Stuffs = newStuffs
	} else {
		if token != "" {
			for _, stuff := range transaction.Stuffs {
				if stuff.Data == nil {
					continue
				}
				decryptedStockValues, err := pkg.Decrypt(stuff.Data.Values, viper.GetString("stock.secret_key"))
				if err != nil {
					logrus.Errorf("failed to decrypt stock values: %v", err)
					return nil, app.NewInternalServerError()
				}
				stuff.Data.Values = decryptedStockValues
			}
		}
	}

	return &transaction, nil
}

func (s *TransactionService) Update(ctx context.Context, transaction *model.Transaction) error {
	err := s.db.Save(&transaction).Error
	if err != nil {
		logrus.Errorf("gorm failed to update transaction: %v", err)
		return app.NewInternalServerError()
	}
	return nil
}
