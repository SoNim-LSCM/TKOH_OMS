package service

import (
	"log"
	"os"
	"time"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/robfig/cron/v3"
)

func SetupCronJob(f *os.File) {
	logCronJob(f)
	log.Println("SYSTEM RESTARTED")

	environment := os.Getenv("ENVIRONMENT")
	if environment == "deployment" {
		c1 := cron.New()
		c1.AddFunc("0 0 * * *", func() {
			logCronJob(f)
			log.Println("START OF A NEW LOG FILE")
			routineCronJob()
		})
		c1.Start()
		defer f.Close()

		for {
			time.Sleep(time.Minute)
		}
	}
}

func logCronJob(f *os.File) {
	logPath := os.Getenv("LOG_PATH")
	now := time.Now()
	f, err := os.OpenFile(logPath+"/TKOH-OMS-LOGS-"+now.Format("2006-01-02"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errorHandler.CheckFatalError(err)
	log.SetOutput(f)
}

func routineCronJob() {
	BackgroundRoutinesToSchedules()
}
