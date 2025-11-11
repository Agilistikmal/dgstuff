package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/paginated"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*fiber.App, *service.StuffService, *gorm.DB) {
	db := app.NewDatabase("sqlite", ":memory:")
	validator := app.NewValidator()
	stuffService := service.NewStuffService(db, validator)
	handler := NewStuffHandler(stuffService)

	app := fiber.New()
	handler.InitRoutes(app)

	assert.NotNil(t, db)
	assert.NotNil(t, validator)
	assert.NotNil(t, stuffService)
	assert.NotNil(t, handler)
	assert.NotNil(t, app)

	return app, stuffService, db
}

func cleanupTest(db *gorm.DB) {
	db.Migrator().DropTable(&model.Stuff{}, &model.StuffCategory{}, &model.StuffMedia{}, "stuff_category_relation")
}

func TestStuffHandler_CreateCategoryNotExists(t *testing.T) {
	app, _, db := setupTest(t)
	defer cleanupTest(db)

	body := bytes.NewBuffer([]byte(`
		{
			"name": "Test Stuff",
			"description": "Test Description",
			"price": 100000,
			"currency": "IDR",
			"stock_count": 100,
			"is_active": true,
			"medias": [
				{
					"url": "https://example.com/image.jpg",
					"type": "image",
					"position": 1
				},
				{
					"url": "https://example.com/video.mp4",
					"type": "video",
					"position": 2
				}
			],
			"categories": ["Account", "Game", "Minecraft"]
		}
	`))

	req := httptest.NewRequest("POST", "/api/stuff", body)
	req.Header.Set("Content-Type", "application/json")

	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var responseBody model.Stuff
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(responseBody.Categories))
	assert.Equal(t, "Account", responseBody.Categories[0].Name)
	assert.Equal(t, "Game", responseBody.Categories[1].Name)
	assert.Equal(t, "Minecraft", responseBody.Categories[2].Name)
}

func TestStuffHandler_CreateSuccess(t *testing.T) {
	app, _, db := setupTest(t)
	defer cleanupTest(db)

	body := bytes.NewBuffer([]byte(`
		{
			"name": "Test Stuff",
			"description": "Test Description",
			"price": 100000,
			"currency": "IDR",
			"stock_count": 100,
			"is_active": true,
			"medias": [
				{
					"url": "https://example.com/image.jpg",
					"type": "image",
					"position": 1
				},
				{
					"url": "https://example.com/video.mp4",
					"type": "video",
					"position": 2
				}
			],
			"categories": []
		}
	`))

	req := httptest.NewRequest("POST", "/api/stuff", body)
	req.Header.Set("Content-Type", "application/json")

	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestStuffHandler_GetBySlugSuccess(t *testing.T) {
	app, stuffService, db := setupTest(t)
	defer cleanupTest(db)

	stuff, err := stuffService.Create(context.Background(), model.StuffCreateDTO{
		Name:        "Test Stuff",
		Description: "Test Description",
		Price:       100000,
		Currency:    "IDR",
		StockCount:  100,
		IsActive:    true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, stuff)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/stuff/%s", stuff.Slug), nil)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.Stuff
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, stuff.ID, responseBody.ID)
	assert.Equal(t, stuff.Name, responseBody.Name)
	assert.Equal(t, stuff.Description, responseBody.Description)
	assert.Equal(t, stuff.Price, responseBody.Price)
	assert.Equal(t, stuff.Currency, responseBody.Currency)
	assert.Equal(t, stuff.StockCount, responseBody.StockCount)
	assert.Equal(t, stuff.IsActive, responseBody.IsActive)
}

func TestStuffHandler_GetAllSuccess(t *testing.T) {
	app, stuffService, db := setupTest(t)
	defer cleanupTest(db)

	stuff, err := stuffService.Create(context.Background(), model.StuffCreateDTO{
		Name:        "Test Stuff",
		Description: "Test Description",
		Price:       100000,
		Currency:    "IDR",
		StockCount:  100,
		IsActive:    true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, stuff)

	req := httptest.NewRequest("GET", "/api/stuff", nil)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody paginated.Paginated[model.Stuff]
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, stuff.ID, responseBody.Data[0].ID)
	assert.Equal(t, stuff.Name, responseBody.Data[0].Name)
	assert.Equal(t, stuff.Description, responseBody.Data[0].Description)
	assert.Equal(t, stuff.Price, responseBody.Data[0].Price)
	assert.Equal(t, stuff.Currency, responseBody.Data[0].Currency)
	assert.Equal(t, stuff.StockCount, responseBody.Data[0].StockCount)
	assert.Equal(t, stuff.IsActive, responseBody.Data[0].IsActive)
}

func TestStuffHandler_GetByCategorySuccess(t *testing.T) {
	app, stuffService, db := setupTest(t)
	defer cleanupTest(db)

	stuff, err := stuffService.Create(context.Background(), model.StuffCreateDTO{
		Name:        "Test Stuff",
		Description: "Test Description",
		Price:       100000,
		Currency:    "IDR",
		StockCount:  100,
		IsActive:    true,
		Categories:  []string{"Account", "Game", "Minecraft"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, stuff)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/stuff/category/%d", stuff.Categories[0].ID), nil)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody paginated.Paginated[model.Stuff]
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, stuff.ID, responseBody.Data[0].ID)
	assert.Equal(t, stuff.Name, responseBody.Data[0].Name)
	assert.Equal(t, stuff.Description, responseBody.Data[0].Description)
	assert.Equal(t, stuff.Price, responseBody.Data[0].Price)
	assert.Equal(t, stuff.Currency, responseBody.Data[0].Currency)
	assert.Equal(t, stuff.StockCount, responseBody.Data[0].StockCount)
	assert.Equal(t, stuff.IsActive, responseBody.Data[0].IsActive)
}
