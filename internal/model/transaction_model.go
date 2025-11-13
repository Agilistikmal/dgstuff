package model

import (
	"time"

	"github.com/agilistikmal/dgstuff/internal/pkg/payment"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	ID        string             `json:"id" gorm:"primaryKey"`
	Email     string             `json:"email"`
	Amount    float64            `json:"amount"`
	Currency  Currency           `json:"currency"`
	Status    TransactionStatus  `json:"status"`
	Stuffs    []TransactionStuff `json:"stuffs" gorm:"foreignKey:TransactionID"`
	Payment   TransactionPayment `json:"payment" gorm:"foreignKey:TransactionID"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TransactionStuff struct {
	ID            int                   `json:"id" gorm:"primaryKey;autoIncrement"`
	TransactionID string                `json:"transaction_id"`
	StuffID       int                   `json:"stuff_id"`
	StuffName     string                `json:"stuff_name"`
	StuffSlug     string                `json:"stuff_slug"`
	StuffPrice    float64               `json:"stuff_price"`
	Quantity      int                   `json:"quantity"`
	TotalPrice    float64               `json:"total_price"`
	Currency      Currency              `json:"currency"`
	Data          *TransactionStuffData `json:"data,omitempty" gorm:"foreignKey:TransactionStuffID"`
}

type TransactionStuffData struct {
	ID                 int    `json:"id"`
	TransactionStuffID int    `json:"transaction_stuff_id"`
	Values             string `json:"values"`
	Separator          string `json:"separator"`
}

type TransactionPayment struct {
	ID            string                  `json:"id" gorm:"primaryKey"`
	TransactionID string                  `json:"transaction_id"`
	Type          payment.PaymentType     `json:"type"`
	Method        payment.PaymentMethod   `json:"method"`
	Code          string                  `json:"code"`
	Status        payment.PaymentStatus   `json:"status"`
	Amount        float64                 `json:"amount"`
	Currency      string                  `json:"currency"`
	Provider      payment.PaymentProvider `json:"provider"`
	URL           string                  `json:"url"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// DTO

type TransactionCreateDTO struct {
	Email           string                      `json:"email" validate:"required,email"`
	Currency        string                      `json:"currency" validate:"required,oneof=USD IDR"`
	Status          string                      `json:"status" validate:"required,oneof=pending success failed"`
	Stuffs          []TransactionStuffCreateDTO `json:"stuffs" validate:"required,max=10"`
	PaymentProvider payment.PaymentProvider     `json:"payment_provider" validate:"required,oneof=xendit midtrans paypal stripe"`
}

type TransactionStuffCreateDTO struct {
	StuffID  int `json:"stuff_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required,min=1"`
}
