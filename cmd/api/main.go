package main

import (
	"fmt"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/agilistikmal/dgstuff/internal/http/handler"
	"github.com/agilistikmal/dgstuff/internal/http/middleware"
	"github.com/agilistikmal/dgstuff/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()

	db := app.NewDatabase(viper.GetString("database.provider"), viper.GetString("database.url"))
	validator := app.NewValidator()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	stuffService := service.NewStuffService(db, validator)
	stuffHandler := handler.NewStuffHandler(stuffService)
	stuffHandler.InitRoutes(app)

	host := "0.0.0.0"
	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	fmt.Printf("--------------------------------\n")
	fmt.Printf("Starting server on %s\n", addr)
	fmt.Printf("Frontend URL: http://%s\n", addr)
	fmt.Printf("Backend API URL: http://%s/api\n", addr)
	fmt.Printf("--------------------------------\n")

	app.Use(middleware.NewInvalidMiddleware().Handle)
	app.Listen(fmt.Sprintf(":%d", port))
}
