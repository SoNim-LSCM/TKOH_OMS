package service

import (
	"strings"
	"time"
)

func StringToResponseTime(timeString string) (string, error) {
	var outputString string
	if timeString == "" {
		timeString = time.Time{}.Format("2006-01-02T15:04:05")
	}
	timeString = strings.Split(timeString, "+")[0]
	timeObj, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		timeObj, err = time.Parse("200601021504", timeString)
		if err != nil {
			timeObj, err = time.Parse("2006-01-02 15:04:05", timeString)
			if err != nil {
				return outputString, err
			}
		}
	}
	return timeObj.Format("200601021504"), nil
}

func StringToDatetime(timeString string) (string, error) {
	var outputString string
	if timeString == "" {
		timeString = time.Time{}.Format("2006-01-02T15:04:05")
	}
	timeString = strings.Split(timeString, "+")[0]
	timeObj, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		timeObj, err = time.Parse("200601021504", timeString)
		if err != nil {
			timeObj, err = time.Parse("2006-01-02 15:04:05", timeString)
			if err != nil {
				return outputString, err
			}
		}
	}
	return timeObj.Format("2006-01-02 15:04:05"), nil
}
