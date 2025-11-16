package model

import (
	"strings"
	"time"

	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/spf13/viper"
)

type Stock struct {
	StuffID   int        `json:"stuff_id" gorm:"primaryKey"`
	Values    string     `json:"values,omitempty"`
	Separator string     `json:"separator,omitempty"`
	Count     int        `json:"count"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (s *Stock) CountValues() int {
	if s.Values == "" {
		return 0
	}

	decryptedValues, err := pkg.Decrypt(s.Values, viper.GetString("stock.secret_key"))
	if err != nil {
		return 0
	}
	return len(strings.Split(decryptedValues, s.Separator))
}

// DTO

type StockUpdateDTO struct {
	Values    string `json:"values" validate:"required"`
	Separator string `json:"separator" validate:"required"`
}
