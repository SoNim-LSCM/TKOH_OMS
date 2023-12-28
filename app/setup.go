package app

import (
	"errors"
	"log"
	"os"

	apiHandler "github.com/SoNim-LSCM/TKOH_OMS/api"
	"github.com/SoNim-LSCM/TKOH_OMS/config"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/service"

	// "github.com/SoNim-LSCM/TKOH_OMS/mqtt"

	"github.com/SoNim-LSCM/TKOH_OMS/router"
	"github.com/SoNim-LSCM/TKOH_OMS/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupAndRunApp() {

	// load env
	err := config.LoadENV()
	errorHandler.CheckError(err, "load env")

	// set output logs
	var f *os.File
	go service.SetupCronJob(f)

	apiHandler.Init()

	db_connected := make(chan bool)

	// start database
	go database.StartMySql(db_connected)
	errorHandler.CheckError(err, "Start MySql")

	// start mqtt server
	// go mqtt.MqttSetup()
	// errorHandler.CheckError(err, "start MQTT")

	// create app
	app := fiber.New()

	// attach middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// setup routes
	router.SetupRoutes(app)

	// attach swagger
	config.AddSwaggerRoutes(app)

	err = errors.New("nil")
	for err != nil {
		err = service.GetLocationFromRFMS()
	}

	// err = service.BackgroundRoutinesToSchedules()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// setup websocket
	go websocket.SetupWebsocket()

	//background service
	go service.BackgroundService(db_connected)

	// get the port and start
	port := os.Getenv("API_PORT")
	app.Listen(":" + port)

	log.Println("FINISH SYSTEM CONFIG")

}
