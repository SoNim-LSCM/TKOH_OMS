package service

import (
	"log"
	"strings"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	"gorm.io/gorm/clause"
)

func FindRecords(records interface{}, table string, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: FindOrders: %s\n", filterFields)
	if err := database.DB.Table(table).Where(filterFields, filterValues...).Find(records).Error; err != nil {
		return err
	}
	// if len(args) == 0 {
	// 	if err := database.DB.Table(table).Where(filter).Find(&records).Error; err != nil {
	// 		return err
	// 	}
	// } else {
	// 	if err := database.DB.Table(table).Where(filter, args...).Find(&records).Error; err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func UpdateRecords(updatedRecordList interface{}, table string, updateMap map[string]interface{}, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	database.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Table(table).Where(filterFields, filterValues...).Updates(updateMap)
	err := FindRecords(updatedRecordList, table, filterFields, filterValues...)
	if err != nil {
		return err
	}
	return nil
}

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
