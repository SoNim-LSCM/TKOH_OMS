package service

import (
	"time"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
)

func BackgroundService(db_connected <-chan bool) {
	backgroundInitOrder(db_connected)
}

func backgroundInitOrder(db_connected <-chan bool) {
	for {
		if <-db_connected {
			BackgroundInitOrderToRFMS()
			err := BackgroundInitOrderToRFMS()
			errorHandler.CheckError(err, "Background process")
			time.Sleep(5 * time.Second)
		}
	}
}
