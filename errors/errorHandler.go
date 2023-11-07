package errorHandler

import (
	"log"
)

func CheckError(err error, action string) bool {
	if err != nil {
		log.Printf("Failure when %s: %s", action, err.Error())
		return true
	}
	return false
}

func CheckFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
