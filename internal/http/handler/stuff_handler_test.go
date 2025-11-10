package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*fiber.App, *service.StuffService) {
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

	return app, stuffService
}

func cleanupTest() {
	db := app.NewDatabase("sqlite", ":memory:")
	db.Migrator().DropTable(&model.Stuff{}, &model.StuffCategory{}, &model.StuffMedia{})
}

func TestStuffHandler_CreateCategoryNotExists(t *testing.T) {
	app, _ := setupTest(t)
	defer cleanupTest()

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
			"categories": [1, 2, 3]
		}
	`))

	req := httptest.NewRequest("POST", "/api/stuff", body)
	req.Header.Set("Content-Type", "application/json")

	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)
	assert.Contains(t, string(responseBody), "some categories not found or invalid")
}

func TestStuffHandler_CreateSuccess(t *testing.T) {
	app, _ := setupTest(t)
	defer cleanupTest()

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
