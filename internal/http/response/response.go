package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *fiber.Ctx, code int, data interface{}) error {
	return c.Status(code).JSON(data)
}

func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(&Response{
		Status:  http.StatusText(code),
		Message: message,
	})
}
