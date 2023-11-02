package app

import (
	"log"
	"os"
	"time"

	"tkoh_oms/config"
	"tkoh_oms/errors"

	// "tkoh_oms/mqtt"
	"tkoh_oms/database"
	"tkoh_oms/router"
	"tkoh_oms/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupAndRunApp() {

	// port := os.Getenv("API_PORT")

	// load env
	err := config.LoadENV()
	errors.CheckError(err, "load env")

	// set output logs
	now := time.Now()
	f, err := os.OpenFile("logs/TKOH-OMS-LOGS-"+now.Format("2006-01-02"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errors.CheckFatalError(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println("SYSTEM RESTARTED")

	// start database
	go database.StartMySql()
	// errors.CheckError(err, "start MySql")

	// start mqtt server
	// go mqtt.MqttSetup()
	// errors.CheckError(err, "start MQTT")

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
