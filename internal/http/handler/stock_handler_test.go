package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type stockTest struct {
	app          *fiber.App
	stuffService *service.StuffService
	stockService *service.StockService
	db           *gorm.DB
}

func setupStockTest(t *testing.T) *stockTest {
	config.LoadConfig()
	db := app.NewDatabase("sqlite", ":memory:")
	validator := app.NewValidator()
	stuffService := service.NewStuffService(db, validator)
	stockService := service.NewStockService(db, validator)
	handler := NewStockHandler(stockService)

	app := fiber.New()
	handler.InitRoutes(app)

	assert.NotNil(t, db)
	assert.NotNil(t, validator)
	assert.NotNil(t, stockService)
	assert.NotNil(t, handler)
	assert.NotNil(t, app)

	return &stockTest{
		app:          app,
		stuffService: stuffService,
		stockService: stockService,
		db:           db,
	}
}

func (s *stockTest) cleanup() {
	s.db.Migrator().DropTable(&model.Stock{})
	os.Remove("config.yml")
}

func TestStockHandler_UpdateSuccess(t *testing.T) {
	s := setupStockTest(t)
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

	body := bytes.NewBuffer([]byte(`
		{
			"values": "user123:password123;user456:password456;user789:password789",
			"separator": ";"
		}
	`))
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/stock/%d", stuff.ID), body)
	req.Header.Set("Content-Type", "application/json")
	response, err := s.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/stock/%d", stuff.ID), nil)
	response, err = s.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.Stock
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, 3, responseBody.Count)
	assert.Equal(t, ";", responseBody.Separator)
}
