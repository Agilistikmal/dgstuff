package payment

import (
	"context"
	"time"
)

type PaymentType string

const (
	PaymentTypeDirect   PaymentType = "direct"
	PaymentTypeRedirect PaymentType = "redirect"
)

type PaymentProvider string

const (
	PaymentProviderXendit   PaymentProvider = "xendit"
	PaymentProviderMidtrans PaymentProvider = "midtrans"
	PaymentProviderPaypal   PaymentProvider = "paypal"
	PaymentProviderStripe   PaymentProvider = "stripe"
)

type PaymentMethod string

const (
	PaymentMethodQRIS PaymentMethod = "qris"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID            string          `json:"id"`
	TransactionID string          `json:"transaction_id"`
	Type          PaymentType     `json:"type"`
	Method        PaymentMethod   `json:"method"`
	Code          string          `json:"code"`
	Status        PaymentStatus   `json:"status"`
	Amount        float64         `json:"amount"`
	Currency      string          `json:"currency"`
	Provider      PaymentProvider `json:"provider"`
	Customer      PaymentCustomer `json:"customer"`
	URL           string          `json:"url"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type PaymentCustomer struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
}

type PaymentInvoiceRequest struct {
	TransactionID string               `json:"transaction_id"`
	Amount        float64              `json:"amount"`
	Currency      string               `json:"currency"`
	Customer      PaymentCustomer      `json:"customer"`
	Items         []PaymentInvoiceItem `json:"items"`
}

type PaymentInvoiceItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type PaymentService interface {
	CreateInvoice(ctx context.Context, request PaymentInvoiceRequest) (*Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
}
