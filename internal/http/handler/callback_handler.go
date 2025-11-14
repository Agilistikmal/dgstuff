package handler

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CallbackHandler struct {
	callbackService *service.CallbackService
}

func NewCallbackHandler(callbackService *service.CallbackService) *CallbackHandler {
	return &CallbackHandler{callbackService: callbackService}
}

func (h *CallbackHandler) InitRoutes(app *fiber.App) {
	apiCallback := app.Group("/api/callback")
	apiCallback.Post("/:provider", h.HandleCallback)
}

func (h *CallbackHandler) HandleCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	payloadToken := c.Get("X-Callback-Token")

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}
	err := h.callbackService.HandleCallback(c.Context(), provider, payload, payloadToken)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, nil)
}
