package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type authTest struct {
	app          *fiber.App
	authService  *service.AuthService
	tokenService *service.TokenService
	db           *gorm.DB
}

func setupAuthTest(t *testing.T) *authTest {
	config.LoadConfig()
	db := app.NewDatabase("sqlite", ":memory:")
	validator := app.NewValidator()
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(db, validator, tokenService)
	handler := NewAuthHandler(authService)

	app := fiber.New()
	handler.InitRoutes(app)

	assert.NotNil(t, db)
	assert.NotNil(t, validator)
	assert.NotNil(t, authService)
	assert.NotNil(t, handler)
	assert.NotNil(t, app)

	return &authTest{
		app:          app,
		authService:  authService,
		tokenService: tokenService,
		db:           db,
	}
}

func (a *authTest) cleanup() {
	a.db.Migrator().DropTable(&model.User{})
}

func TestAuthHandler_RegisterSuccess(t *testing.T) {
	a := setupAuthTest(t)
	defer a.cleanup()

	body := bytes.NewBuffer([]byte(`
		{
			"email": "test@example.com",
			"password": "password",
			"confirm_password": "password"
		}
	`))
	req := httptest.NewRequest("POST", "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	response, err := a.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.AuthResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.NotNil(t, responseBody)
	assert.NotEmpty(t, responseBody.Token)
	t.Logf("token: %s", responseBody.Token)
}

func TestAuthHandler_LoginSuccess(t *testing.T) {
	a := setupAuthTest(t)
	defer a.cleanup()

	auth, err := a.authService.Register(context.Background(), model.UserRegisterDTO{
		Email:           "test@example.com",
		Password:        "password",
		ConfirmPassword: "password",
	})
	assert.NoError(t, err)
	assert.NotNil(t, auth)

	body := bytes.NewBuffer([]byte(`
		{
			"email": "test@example.com",
			"password": "password"
		}
	`))
	req := httptest.NewRequest("POST", "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	response, err := a.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.AuthResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.NotNil(t, responseBody)
	assert.NotEmpty(t, responseBody.Token)
	t.Logf("token: %s", responseBody.Token)
}
