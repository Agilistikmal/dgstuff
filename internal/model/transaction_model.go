package model

import "time"

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
	Currency  string             `json:"currency"`
	Status    TransactionStatus  `json:"status"`
	Stuffs    []TransactionStuff `json:"stuffs" gorm:"foreignKey:TransactionID"`
	Payment   TransactionPayment `json:"payment" gorm:"foreignKey:TransactionID"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TransactionStuff struct {
	ID            int      `json:"id" gorm:"primaryKey;autoIncrement"`
	TransactionID string   `json:"transaction_id"`
	StuffID       int      `json:"stuff_id"`
	StuffName     string   `json:"stuff_name"`
	StuffSlug     string   `json:"stuff_slug"`
	StuffPrice    float64  `json:"stuff_price"`
	Quantity      int      `json:"quantity"`
	TotalPrice    float64  `json:"total_price"`
	Currency      Currency `json:"currency"`
}

type TransactionPayment struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	TransactionID    string    `json:"transaction_id"`
	PaymentMethod    string    `json:"payment_method"`
	PaymentCode      string    `json:"payment_code"`
	PaymentStatus    string    `json:"payment_status"`
	PaymentAmount    float64   `json:"payment_amount"`
	PaymentCurrency  Currency  `json:"payment_currency"`
	PaymentProvider  string    `json:"payment_provider"`
	PaymentReference string    `json:"payment_reference"`
	PayerName        string    `json:"payer_name"`
	PaymentCreatedAt time.Time `json:"payment_created_at"`
	PaymentUpdatedAt time.Time `json:"payment_updated_at"`
}
