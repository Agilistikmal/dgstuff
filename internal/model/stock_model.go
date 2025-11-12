package model

import (
	"strings"
	"time"
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
	return len(strings.Split(s.Values, s.Separator))
}

// DTO

type StockUpdateDTO struct {
	Values    string `json:"values" validate:"required"`
	Separator string `json:"separator" validate:"required"`
}
