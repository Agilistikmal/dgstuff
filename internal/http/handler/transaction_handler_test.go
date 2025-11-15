package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment/xendit_payment"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type transactionTest struct {
	app                *fiber.App
	xenditPayment      payment.PaymentService
	stuffService       *service.StuffService
	stockService       *service.StockService
	transactionService *service.TransactionService
	db                 *gorm.DB
}

func setupTransactionTest(t *testing.T) *transactionTest {
	config.LoadConfig()
	db := app.NewDatabase("sqlite", ":memory:")
	validator := app.NewValidator()
	xenditPayment := xendit_payment.NewXenditPayment(viper.GetString("payment.provider.xendit.api_key"))
	stuffService := service.NewStuffService(db, validator)
	stockService := service.NewStockService(db, validator)
	transactionService := service.NewTransactionService(db, validator, xenditPayment)
	handler := NewTransactionHandler(transactionService)

	app := fiber.New()
	handler.InitRoutes(app)

	assert.NotNil(t, db)
	assert.NotNil(t, validator)
	assert.NotNil(t, xenditPayment)
	assert.NotNil(t, transactionService)
	assert.NotNil(t, handler)
	assert.NotNil(t, app)

	return &transactionTest{
		app:                app,
		xenditPayment:      xenditPayment,
		stuffService:       stuffService,
		stockService:       stockService,
		transactionService: transactionService,
		db:                 db,
	}
}

func (t *transactionTest) cleanup() {
	t.db.Migrator().DropTable(&model.Transaction{}, &model.TransactionStuff{}, &model.TransactionPayment{})
}

func TestTransactionHandler_CreateSuccess(t *testing.T) {
	s := setupTransactionTest(t)
	defer s.cleanup()

	stuff, err := s.stuffService.Create(context.Background(), model.StuffCreateDTO{
		Name:        "Test Stuff",
		Description: "Test Description",
		Price:       100000,
		Currency:    "IDR",
		IsActive:    true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, stuff)

	err = s.stockService.Update(context.Background(), stuff.ID, model.StockUpdateDTO{
		Values:    "user123:password123;user456:password456;user789:password789",
		Separator: ";",
	})
	assert.NoError(t, err)

	body := bytes.NewBuffer([]byte(`
		{
			"email": "test@example.com",
			"currency": "IDR",
			"stuffs": [
				{
					"stuff_id": ` + strconv.Itoa(stuff.ID) + `,
					"quantity": 1
				}
			],
			"payment_provider": "xendit"
		}
	`))
	req := httptest.NewRequest("POST", "/api/transaction", body)
	req.Header.Set("Content-Type", "application/json")
	response, err := s.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var transaction model.Transaction
	err = json.NewDecoder(response.Body).Decode(&transaction)
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, model.TransactionStatusPending, transaction.Status)
	assert.Equal(t, "test@example.com", transaction.Email)
	assert.Equal(t, model.CurrencyIDR, transaction.Currency)
	assert.Equal(t, 1, len(transaction.Stuffs))
	assert.Equal(t, stuff.ID, transaction.Stuffs[0].StuffID)
	t.Logf("payment url: %s", transaction.Payment.URL)
}
