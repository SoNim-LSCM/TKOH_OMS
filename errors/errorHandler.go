package errors

import (
	"log"
)

func CheckError(err error, action string) {
	if err != nil {
		log.Printf("Failure when %s: %s", action, err.Error())
		panic(err)
	}
}

func CheckFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
