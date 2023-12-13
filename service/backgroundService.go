package service

import (
	"time"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
)

func BackgroundService(db_connected <-chan bool) {
	for {
		if <-db_connected {
			go backgroundInitOrder()
			return
		}
	}
}

func backgroundInitOrder() {
	for {
		// BackgroundInitOrderToRFMS()
		err := BackgroundInitOrderToRFMS()
		errorHandler.CheckError(err, "Background process")
		time.Sleep(5 * time.Second)
	}
}
