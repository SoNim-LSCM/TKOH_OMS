package handlers

import (
	"tkoh_oms/service"

	"github.com/gofiber/fiber/v2"
)

// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Health
// @Accept */*
// @Produce plain
// @Success 200 "OK"
// @Router /health [get]
func HandleHealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

// @Summary Toggle Background Service.
// @Description Toggle Background Service
// @Tags Health
// @Accept */*
// @Produce plain
// @Success 200 "Background service is ON / OFF"
// @Router /toggleBackgroundService [get]
func ToggleBackgroundService(c *fiber.Ctx) error {
	if service.ToggleBackgroundInitOrder() {
		return c.SendString("Background service is ON")
	}
	return c.SendString("Background service is OFF")
}
