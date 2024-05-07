package config

import (
	"os"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "tkoh_oms/docs"
)

func AddSwaggerRoutes(app *fiber.App) {
	// setup swagger
	url := os.Getenv("URL")
	app.Get("/"+url+"/swagger/*", fiberSwagger.WrapHandler)
}
