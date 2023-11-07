package app

import (
	"log"
	"os"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/config"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"

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
	now := time.Now()
	logPath := os.Getenv("LOG_PATH")
	f, err := os.OpenFile(logPath+"/TKOH-OMS-LOGS-"+now.Format("2006-01-02"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errorHandler.CheckFatalError(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println("SYSTEM RESTARTED")

	// start database
	go database.StartMySql()
	errorHandler.CheckError(err, "start MySql")

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

	// setup websocket
	go websocket.SetupWebsocket()

	// get the port and start
	port := os.Getenv("API_PORT")
	app.Listen(":" + port)

	log.Println("FINISH SYSTEM CONFIG")
}
