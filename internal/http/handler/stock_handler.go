package handler

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
)

type StockHandler struct {
	stockService *service.StockService
}

func NewStockHandler(stockService *service.StockService) *StockHandler {
	return &StockHandler{stockService: stockService}
}

func (h *StockHandler) InitRoutes(app *fiber.App) {
	apiStock := app.Group("/api/stock")
	apiStock.Put("/:stuff_id", h.Update)
	apiStock.Get("/:stuff_id", h.Get)
}

func (h *StockHandler) Update(c *fiber.Ctx) error {
	stuffID, err := c.ParamsInt("stuff_id")
	if err != nil {
		return response.Error(c, app.NewBadRequestError("stuff id must be an integer"))
	}
	var dto model.StockUpdateDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}
	err = h.stockService.Update(c.Context(), stuffID, dto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, nil)
}

func (h *StockHandler) Get(c *fiber.Ctx) error {
	stuffID, err := c.ParamsInt("stuff_id")
	if err != nil {
		return response.Error(c, app.NewBadRequestError("stuff id must be an integer"))
	}
	stock, err := h.stockService.Get(c.Context(), stuffID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, stock)
}
