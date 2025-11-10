package response

import (
	"errors"
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/app"
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

func Error(c *fiber.Ctx, err error) error {
	var appErr *app.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Code).JSON(&Response{
			Status:  http.StatusText(appErr.Code),
			Message: appErr.Message,
		})
	} else {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "internal server error",
		})
	}
}
