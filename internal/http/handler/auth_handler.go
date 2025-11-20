package handler

import (
	"net/http"
	"time"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) InitRoutes(app *fiber.App) {
	apiAuth := app.Group("/api/auth")
	apiAuth.Post("/login", h.Login)
	apiAuth.Post("/register", h.Register)
	apiAuth.Get("/me", h.Me)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	token := c.Cookies("auth_token", "")
	user, err := h.authService.Me(c.Context(), token)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var dto model.UserLoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}
	auth, err := h.authService.Login(c.Context(), dto)
	if err != nil {
		return response.Error(c, err)
	}
	c.Cookie(&fiber.Cookie{
		Name:    "auth_token",
		Value:   auth.Token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return response.Success(c, http.StatusOK, auth)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var dto model.UserRegisterDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}
	auth, err := h.authService.Register(c.Context(), dto)
	if err != nil {
		return response.Error(c, err)
	}
	c.Cookie(&fiber.Cookie{
		Name:    "auth_token",
		Value:   auth.Token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return response.Success(c, http.StatusOK, auth)
}
