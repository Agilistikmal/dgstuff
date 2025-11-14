package service

import (
	"context"
	"time"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xendit/xendit-go/v7/invoice"
)

type CallbackService struct {
	transactionService *TransactionService
}

func NewCallbackService(transactionService *TransactionService) *CallbackService {
	return &CallbackService{transactionService: transactionService}
}

func (s *CallbackService) HandleCallback(ctx context.Context, provider string, payload map[string]interface{}, payloadToken string) error {
	if payloadToken == "" {
		logrus.Errorf("callback token is required")
		return app.NewForbiddenError("callback token is required")
	}

	providerEnum := payment.PaymentProvider(provider)

	switch providerEnum {
	case payment.PaymentProviderXendit:
		return s.handleXenditCallback(ctx, payload, payloadToken)
	default:
		logrus.Errorf("payment provider not supported: %s", provider)
		return app.NewBadRequestError("payment provider not supported")
	}
}

func (s *CallbackService) handleXenditCallback(ctx context.Context, payload map[string]interface{}, payloadToken string) error {
	logrus.Info("Trying to handle Xendit callback")
	invoiceCallbackJson, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("failed to marshal invoice callback: %v", err)
		return app.NewInternalServerError()
	}
	invoiceCallback := invoice.InvoiceCallback{}
	err = json.Unmarshal(invoiceCallbackJson, &invoiceCallback)
	if err != nil {
		logrus.Errorf("failed to unmarshal invoice callback: %v", err)
		return app.NewBadRequestError("invalid invoice callback")
	}

	callbackToken := viper.GetString("payment.provider.xendit.webhook_token")
	if payloadToken != callbackToken {
		logrus.Errorf("invalid callback token: %s", payloadToken)
		return app.NewForbiddenError("invalid callback token")
	}

	transaction, err := s.transactionService.Get(ctx, invoiceCallback.ExternalId)
	if err != nil {
		logrus.Errorf("failed to get transaction: %v", err)
		return err
	}

	if transaction.Payment.ID != invoiceCallback.Id {
		logrus.Errorf("transaction id does not match: %s != %s", transaction.Payment.ID, invoiceCallback.Id)
		return app.NewBadRequestError("transaction id does not match")
	}

	switch invoiceCallback.Status {
	case "PAID":
		transaction.Status = model.TransactionStatusSuccess
	case "EXPIRED":
		transaction.Status = model.TransactionStatusFailed
	case "FAILED":
		transaction.Status = model.TransactionStatusFailed
	default:
		logrus.Errorf("invalid invoice callback status: %s", invoiceCallback.Status)
		return app.NewBadRequestError("invalid invoice callback status")
	}
	transaction.Payment.Status = payment.PaymentStatus(invoiceCallback.Status)
	transaction.Payment.UpdatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	err = s.transactionService.Update(ctx, transaction)
	if err != nil {
		logrus.Errorf("failed to update transaction: %v", err)
		return app.NewInternalServerError()
	}

	return nil
}
