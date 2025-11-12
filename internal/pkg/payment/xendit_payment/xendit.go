package xendit_payment

import (
	"context"

	"github.com/agilistikmal/dgstuff/internal/pkg/payment"
	"github.com/spf13/viper"
	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type XenditPayment struct {
	client *xendit.APIClient
}

func NewXenditPayment(apiKey string) payment.PaymentService {
	client := xendit.NewClient(apiKey)
	return &XenditPayment{client: client}
}

func (p *XenditPayment) CreateInvoice(ctx context.Context, request payment.PaymentInvoiceRequest) (*payment.Payment, error) {
	invReq := invoice.CreateInvoiceRequest{
		ExternalId: request.TransactionID,
		Amount:     request.Amount,
		Currency:   &request.Currency,
		Customer: &invoice.CustomerObject{
			Email:       *invoice.NewNullableString(&request.Customer.Email),
			GivenNames:  *invoice.NewNullableString(&request.Customer.Name),
			PhoneNumber: *invoice.NewNullableString(&request.Customer.Phone),
		},
	}
	expirationTime := viper.GetDuration("payment.expiration_time")
	expirationTimeFloat := float32(expirationTime.Seconds())
	invReq.InvoiceDuration = &expirationTimeFloat

	for _, item := range request.Items {
		invReq.Items = append(invReq.Items, invoice.InvoiceItem{
			Name:        item.Name,
			Quantity:    float32(item.Quantity),
			Price:       float32(item.Price),
			ReferenceId: &item.ID,
			Category:    &item.ID,
		})
	}

	inv, _, err := p.client.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(invReq).
		Execute()
	if err != nil {
		return nil, err
	}

	payment := payment.Payment{
		ID:            *inv.Id,
		TransactionID: inv.ExternalId,
		Amount:        inv.Amount,
		Currency:      string(*inv.Currency),
		Provider:      payment.PaymentProviderXendit,
		Customer: payment.PaymentCustomer{
			Name:  *inv.Customer.GivenNames.Get(),
			Email: *inv.Customer.Email.Get(),
			Phone: *inv.Customer.PhoneNumber.Get(),
		},
		Status:    payment.PaymentStatusPending,
		CreatedAt: inv.Created,
		UpdatedAt: inv.Updated,
		Type:      payment.PaymentTypeRedirect,
		URL:       inv.InvoiceUrl,
	}

	return &payment, nil
}

func (p *XenditPayment) GetPayment(ctx context.Context, paymentID string) (*payment.Payment, error) {
	inv, _, err := p.client.InvoiceApi.GetInvoiceById(ctx, paymentID).
		Execute()
	if err != nil {
		return nil, err
	}

	return &payment.Payment{
		ID:            *inv.Id,
		TransactionID: inv.ExternalId,
		Amount:        inv.Amount,
		Currency:      string(*inv.Currency),
		Provider:      payment.PaymentProviderXendit,
		Customer: payment.PaymentCustomer{
			Name:  *inv.Customer.GivenNames.Get(),
			Email: *inv.Customer.Email.Get(),
			Phone: *inv.Customer.PhoneNumber.Get(),
		},
		Status:    payment.PaymentStatus(inv.Status),
		CreatedAt: inv.Created,
		UpdatedAt: inv.Updated,
		Type:      payment.PaymentTypeRedirect,
		URL:       inv.InvoiceUrl,
	}, nil
}
