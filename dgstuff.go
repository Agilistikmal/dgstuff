package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/agilistikmal/dgstuff/cmd/api"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed ui/dist
var uiDist embed.FS

func main() {
	config.LoadConfig()

	app := api.Run()

	distFS, _ := fs.Sub(uiDist, "ui/dist")
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(distFS),
		Browse: true,
		Index:  "index.html",
	}))

	app.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api") {
			return c.Next()
		}
		index, err := distFS.Open("index.html")
		if err != nil {
			return fiber.ErrNotFound
		}
		defer index.Close()
		data, _ := io.ReadAll(index)
		c.Type("html")
		return c.Send(data)
	})

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

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		logrus.Fatalf("failed to start server: %v", err)
	}
}
