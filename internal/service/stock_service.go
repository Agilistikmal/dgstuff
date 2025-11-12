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

type StockService struct {
	db        *gorm.DB
	validator *app.Validator
}

func NewStockService(db *gorm.DB, validator *app.Validator) *StockService {
	return &StockService{db: db, validator: validator}
}

func (s *StockService) encrypt(value string) (string, error) {
	return pkg.Encrypt(value, viper.GetString("stock.secret_key"))
}

func (s *StockService) decrypt(value string) (string, error) {
	return pkg.Decrypt(value, viper.GetString("stock.secret_key"))
}

func (s *StockService) Update(ctx context.Context, stuffID int, dto model.StockUpdateDTO) error {
	err := s.validator.Validate(dto)
	if err != nil {
		return err
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	encryptedValues, err := s.encrypt(dto.Values)
	if err != nil {
		return err
	}

	stock := model.Stock{
		StuffID:   stuffID,
		Values:    encryptedValues,
		Separator: dto.Separator,
	}

	err = tx.Save(&stock).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to update stock: %v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (s *StockService) Get(ctx context.Context, stuffID int) (*model.Stock, error) {
	var stock model.Stock
	err := s.db.Where("stuff_id = ?", stuffID).First(&stock).Error
	if err != nil {
		return nil, err
	}

	decryptedValues, err := s.decrypt(stock.Values)
	if err != nil {
		return nil, err
	}

	stock.Values = decryptedValues
	stock.Count = stock.CountValues()
	return &stock, nil
}
