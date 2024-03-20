package service

import (
	"log"
	"os"
	"time"

	errorHandler "tkoh_oms/errors"
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
	floorPlan, _ := GetFloorPlan()
	for {
		log.Println("Background Process Running ...")
		if isLooping {
			err := BackgroundInitOrderToRFMS()
			errorHandler.CheckError(err, "Background Init Order to RFMS")
			err = BackgroundReportRobotStatus(floorPlan)
			errorHandler.CheckError(err, "Background Report Robot Status")
		}
		time.Sleep(5 * time.Second)
	}
}
