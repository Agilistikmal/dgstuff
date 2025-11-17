package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment/xendit_payment"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type callbackTest struct {
	app                *fiber.App
	callbackService    *service.CallbackService
	transactionService *service.TransactionService
	stuffService       *service.StuffService
	stockService       *service.StockService
	db                 *gorm.DB
}

func setupCallbackTest(t *testing.T) *callbackTest {
	config.LoadConfig()
	db := app.NewDatabase("sqlite", ":memory:")
	validator := app.NewValidator()
	xenditPayment := xendit_payment.NewXenditPayment(viper.GetString("payment.provider.xendit.api_key"))
	tokenService := service.NewTokenService()
	transactionService := service.NewTransactionService(db, validator, xenditPayment, tokenService)
	callbackService := service.NewCallbackService(transactionService)
	stuffService := service.NewStuffService(db, validator)
	stockService := service.NewStockService(db, validator)
	handler := NewCallbackHandler(callbackService)

	app := fiber.New()
	handler.InitRoutes(app)

	assert.NotNil(t, db)
	assert.NotNil(t, validator)
	assert.NotNil(t, transactionService)
	assert.NotNil(t, callbackService)
	assert.NotNil(t, stuffService)
	assert.NotNil(t, stockService)
	assert.NotNil(t, handler)
	assert.NotNil(t, app)

	return &callbackTest{
		app:                app,
		callbackService:    callbackService,
		transactionService: transactionService,
		stuffService:       stuffService,
		stockService:       stockService,
		db:                 db,
	}
}

func (c *callbackTest) cleanup() {
	c.db.Migrator().DropTable(&model.Transaction{}, &model.TransactionStuff{}, &model.TransactionPayment{}, &model.Stuff{}, &model.Stock{})
	sqlDb, err := c.db.DB()
	if err != nil {
		logrus.Errorf("failed to get sql db: %v", err)
		return
	}
	sqlDb.Close()
}

func TestCallbackHandler_HandleCallbackSuccess(t *testing.T) {
	c := setupCallbackTest(t)
	defer c.cleanup()

	stuff, err := c.stuffService.Create(context.Background(), model.StuffCreateDTO{
		Name:        "Test Stuff",
		Price:       100000,
		Currency:    "IDR",
		Description: "Test Description",
		IsActive:    true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, stuff)

	err = c.stockService.Update(context.Background(), stuff.ID, model.StockUpdateDTO{
		Values:    "user123:password123;user456:password456;user789:password789",
		Separator: ";",
	})
	assert.NoError(t, err)

	transactionDTO := model.TransactionCreateDTO{
		Email:    "test@example.com",
		Currency: "IDR",
		Stuffs: []model.TransactionStuffCreateDTO{
			{
				StuffID:  stuff.ID,
				Quantity: 1,
			},
		},
		PaymentProvider: "xendit",
	}

	transaction, err := c.transactionService.Create(context.Background(), transactionDTO)
	assert.NoError(t, err)
	assert.NotNil(t, transaction)

	xenditInvoiceCallback := map[string]interface{}{
		"id":                  transaction.Payment.ID,
		"external_id":         transaction.ID,
		"user_id":             "5848fdf860053555135587e7",
		"payment_method":      "RETAIL_OUTLET",
		"status":              "PAID",
		"merchant_name":       "Xendit",
		"amount":              transaction.Amount,
		"paid_amount":         transaction.Amount,
		"paid_at":             "2020-01-14T02:32:50.912Z",
		"payer_email":         "test@xendit.co",
		"description":         "Invoice webhook test",
		"created":             "2020-01-13T02:32:49.827Z",
		"updated":             "2020-01-13T02:32:50.912Z",
		"currency":            "IDR",
		"payment_channel":     "ALFAMART",
		"payment_destination": "TEST815",
	}

	xenditInvoiceCallbackJson, err := json.Marshal(xenditInvoiceCallback)
	assert.NoError(t, err)
	body := bytes.NewBuffer(xenditInvoiceCallbackJson)
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/callback/%s", transaction.Payment.Provider), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Callback-Token", viper.GetString("payment.provider.xendit.webhook_token"))
	response, err := c.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
