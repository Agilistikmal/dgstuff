package handler

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) InitRoutes(app *fiber.App) {
	apiTransaction := app.Group("/api/transaction")
	apiTransaction.Post("/", h.Create)
	apiTransaction.Get("/:id", h.Get)
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var dto model.TransactionCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}

	transaction, err := h.transactionService.Create(c.Context(), dto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusCreated, transaction)
}

func (h *TransactionHandler) Get(c *fiber.Ctx) error {
	transactionID := c.Params("id")
	transaction, err := h.transactionService.Get(c.Context(), transactionID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, transaction)
}
