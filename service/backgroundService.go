package service

import (
	"os"
	"time"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
)

var isLooping = true

func BackgroundService(db_connected <-chan bool) {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "deployment" {
		for {
			if <-db_connected {
				go backgroundInitOrder()
				return
			}
		}
	}
}

func ToggleBackgroundInitOrder() bool {
	isLooping = !isLooping
	return isLooping
}

func backgroundInitOrder() {
	for {
		if isLooping {
			err := BackgroundInitOrderToRFMS()
			errorHandler.CheckError(err, "Background Init Order to RFMS")
			err = BackgroundReportRobotStatus()
			errorHandler.CheckError(err, "Background Report Robot Status")
		}
		time.Sleep(5 * time.Second)
	}
}
