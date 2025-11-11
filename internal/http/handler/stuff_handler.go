package handler

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/app"
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
	apiStuff.Get("/:slug", h.GetBySlug)
	apiStuff.Get("/", h.GetAll)
	apiStuff.Get("/category/:category_id", h.GetByCategory)
}

func (h *StuffHandler) Create(c *fiber.Ctx) error {
	var dto model.StuffCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}

	stuff, err := h.stuffService.Create(c.Context(), dto)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, stuff)
}

func (h *StuffHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	stuff, err := h.stuffService.GetBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, stuff)
}

func (h *StuffHandler) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	stuffs, err := h.stuffService.GetAll(c.Context(), page, limit)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, stuffs)
}

func (h *StuffHandler) GetByCategory(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("category_id")
	if err != nil {
		return response.Error(c, app.NewBadRequestError("category id must be an integer"))
	}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	stuffs, err := h.stuffService.GetByCategory(c.Context(), categoryID, page, limit)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, stuffs)
}
