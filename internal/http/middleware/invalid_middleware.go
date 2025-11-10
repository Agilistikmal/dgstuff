package middleware

import (
	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

type InvalidMiddleware struct {
}

func NewInvalidMiddleware() *InvalidMiddleware {
	return &InvalidMiddleware{}
}

func (m *InvalidMiddleware) Handle(c *fiber.Ctx) error {
	return response.Error(c, app.NewBadRequestError("invalid request"))
}
