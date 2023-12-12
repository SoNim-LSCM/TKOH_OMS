package app

import (
	"log"
	"os"
	"time"

	apiHandler "github.com/SoNim-LSCM/TKOH_OMS/api"
	"github.com/SoNim-LSCM/TKOH_OMS/config"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/service"
	"github.com/robfig/cron/v3"

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
	go config.SetupLogCron(f)

	logPath := os.Getenv("LOG_PATH")
	now := time.Now()
	f, err = os.OpenFile(logPath+"/TKOH-OMS-LOGS-"+now.Format("2006-01-02"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errorHandler.CheckFatalError(err)
	c1 := cron.New()
	c1.AddFunc("0 0 * * *", func() {
		// fmt.Println("test")
		f.Close()
		now := time.Now()
		f, err = os.OpenFile(logPath+"/TKOH-OMS-LOGS-"+now.Format("2006-01-02"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		errorHandler.CheckFatalError(err)

		log.SetOutput(f)
		log.Println("START OF A NEW LOG FILE !!!")
	})
	c1.Start()
	log.SetOutput(f)

	defer f.Close()

	log.SetOutput(f)
	log.Println("SYSTEM RESTARTED")

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

	// setup websocket
	go websocket.SetupWebsocket()

	//background service
	go service.BackgroundService(db_connected)

	// get the port and start
	port := os.Getenv("API_PORT")
	app.Listen(":" + port)

	log.Println("FINISH SYSTEM CONFIG")
}
