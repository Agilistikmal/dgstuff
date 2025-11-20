package api

import (
	"strings"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/handler"
	"github.com/agilistikmal/dgstuff/internal/pkg/payment/xendit_payment"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func Run() *fiber.App {
	db := app.NewDatabase(viper.GetString("database.provider"), viper.GetString("database.url"))
	validator := app.NewValidator()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	allowedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:8080",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:8080",
		viper.GetString("app.website_url"),
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, X-Transaction-Token",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	xenditPayment := xendit_payment.NewXenditPayment(viper.GetString("payment.provider.xendit.api_key"))
	tokenService := service.NewTokenService()

	authService := service.NewAuthService(db, validator, tokenService)
	authHandler := handler.NewAuthHandler(authService)
	authHandler.InitRoutes(app)

	appInfoService := service.NewAppInfoService(db)
	appInfoHandler := handler.NewAppInfoHandler(appInfoService)
	appInfoHandler.InitRoutes(app)

	stuffService := service.NewStuffService(db, validator)
	stuffHandler := handler.NewStuffHandler(stuffService)
	stuffHandler.InitRoutes(app)

	stockService := service.NewStockService(db, validator)
	stockHandler := handler.NewStockHandler(stockService)
	stockHandler.InitRoutes(app)

	transactionService := service.NewTransactionService(db, validator, xenditPayment, tokenService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	transactionHandler.InitRoutes(app)

	callbackService := service.NewCallbackService(transactionService)
	callbackHandler := handler.NewCallbackHandler(callbackService)
	callbackHandler.InitRoutes(app)

	return app
}
