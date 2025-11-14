package handler

import (
	"net/http"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/response"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AppInfoHandler struct {
	appInfoService *service.AppInfoService
}

func NewAppInfoHandler(appInfoService *service.AppInfoService) *AppInfoHandler {
	return &AppInfoHandler{appInfoService: appInfoService}
}

func (h *AppInfoHandler) InitRoutes(app *fiber.App) {
	apiAppInfo := app.Group("/api/appinfo")
	apiAppInfo.Get("/", h.Get)
	apiAppInfo.Put("/", h.Update)
}

func (h *AppInfoHandler) Get(c *fiber.Ctx) error {
	appInfo, err := h.appInfoService.Get()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, appInfo)
}

func (h *AppInfoHandler) Update(c *fiber.Ctx) error {
	var appInfo model.AppInfo
	if err := c.BodyParser(&appInfo); err != nil {
		return response.Error(c, app.NewBadRequestError(err.Error()))
	}
	err := h.appInfoService.Update(&appInfo)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, nil)
}
