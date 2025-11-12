package service

import (
	"context"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TransactionService struct {
	db        *gorm.DB
	validator *app.Validator
}

func NewTransactionService(db *gorm.DB, validator *app.Validator) *TransactionService {
	return &TransactionService{db: db, validator: validator}
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
		Amount:   dto.Amount,
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
	}

	err = tx.Save(&transaction).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to create transaction: %v", err)
		return nil, app.NewInternalServerError()
	}

	return &transaction, nil
}
