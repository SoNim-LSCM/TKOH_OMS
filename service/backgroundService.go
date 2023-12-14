package service

import (
	"time"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
)

var isLooping = true

func BackgroundService(db_connected <-chan bool) {
	for {
		if <-db_connected {
			go backgroundInitOrder()
			return
		}
	}
}

func ToggleBackgroundInitOrder() bool {
	isLooping = !isLooping
	return isLooping
}

func backgroundInitOrder() {
	for {
		// BackgroundInitOrderToRFMS()
		if isLooping {
			err := BackgroundInitOrderToRFMS()
			errorHandler.CheckError(err, "Background process")
		}
		time.Sleep(5 * time.Second)
	}
}
