package config

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "github.com/SoNim-LSCM/TKOH_OMS/docs"
)

func AddSwaggerRoutes(app *fiber.App) {
	// setup swagger
	app.Get("/oms/swagger/*", fiberSwagger.WrapHandler)
}
