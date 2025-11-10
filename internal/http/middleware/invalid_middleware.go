package middleware

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

type InvalidMiddleware struct {
}

func NewInvalidMiddleware() *InvalidMiddleware {
	return &InvalidMiddleware{}
}

func (m *InvalidMiddleware) Handle(c *fiber.Ctx) error {
	return response.Error(c, http.StatusBadRequest, "invalid request")
}
