package handler

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
)

type StuffHandler struct {
	stuffService *service.StuffService
}

func NewStuffHandler(stuffService *service.StuffService) *StuffHandler {
	return &StuffHandler{stuffService: stuffService}
}

func (h *StuffHandler) InitRoutes(app *fiber.App) {
	apiStuff := app.Group("/api/stuff")
	apiStuff.Post("/", h.Create)
}

func (h *StuffHandler) Create(c *fiber.Ctx) error {
	var dto model.StuffCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	stuff, err := h.stuffService.Create(c.Context(), dto)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, stuff)
}
