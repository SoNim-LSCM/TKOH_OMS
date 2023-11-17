package router

import (
	"github.com/SoNim-LSCM/TKOH_OMS/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// setup the oms group
	oms := app.Group("/oms")

	// Health Check
	app.Get("/health", handlers.HandleHealthCheck)

	// Login Auth
	oms.Post("/loginStaff", handlers.HandleLoginStaff)
	oms.Post("/loginAdmin", handlers.HandleLoginAdmin)
	oms.Get("/logout", handlers.HandleLogout)
	oms.Get("/renewToken", handlers.HandleRenewToken)

	// Order Management
	oms.Get("/getDeliveryOrder", handlers.HandleGetDeliveryOrder)
	oms.Post("/addDeliveryOrder", handlers.HandleAddDeliveryOrder)
	oms.Get("/getRoutineDeliveryOrder", handlers.HandleGetRoutineDeliveryOrder)
	oms.Get("/triggerHandlingOrder", handlers.HandleTriggerHandlingOrder)
	oms.Post("/updateDeliveryOrder", handlers.HandleUpdateDeliveryOrder)
	oms.Post("/cancelDeliveryOrder", handlers.HandleCancelDeliveryOrder)
	oms.Post("/addRoutine", handlers.HandleAddRoutine)
	oms.Post("/updateRoutineDeliveryOrder", handlers.HandleUpdateRoutineDeliveryOrder)

	// for rfms
	oms.Post("/reportJobStatus", handlers.HandleReportJobStatus)
	oms.Post("/reportSystemStatus", handlers.HandleReportSystemStatus)

	// Map Handling
	oms.Get("/getFloorPlan", handlers.HandleGetFloorPlan)
	oms.Get("/getDutyRooms", handlers.HandleGetDutyRooms)

	// User Management
	oms.Post("/testAW2", handlers.HandleTestAW2)
	oms.Post("/testOW1", handlers.HandleTestOW1)
	oms.Post("/testMW1", handlers.HandleTestMW1)
	oms.Post("/testSW1", handlers.HandleTestSW1)

}
