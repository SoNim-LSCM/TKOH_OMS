package service

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindRecords(db *gorm.DB, records interface{}, table string, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: FindRecords: %s\n", filterFields)
	if err := db.Table(table).Where(filterFields, filterValues...).Find(records).Error; err != nil {
		return err
	}
	return nil
}

func UpdateRecords(db *gorm.DB, updatedRecordList interface{}, table string, updateMap map[string]interface{}, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: UpdateRecords\n")
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Table(table).Where(filterFields, filterValues...).Updates(updateMap)
	err := FindRecords(db, updatedRecordList, table, filterFields, filterValues...)
	if err != nil {
		return err
	}
	return nil
}

func AddRecords(db *gorm.DB, records interface{}) error {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: AddRecords\n")
	if err := db.Create(records).Error; err != nil {
		return err
	}
	return nil
}

func AddSchedulesLogs(db *gorm.DB, userId int, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	logRecord := []db_models.SchedulesLogs{}
	schedules := []db_models.Schedules{}
	err := FindRecords(db, &schedules, "schedules", filterFields, filterValues...)
	if err != nil {
		return err
	}
	logRecord, err = SchedulesToSchedulesLogs(userId, schedules)
	if err != nil {
		return err
	}
	return AddRecords(db, logRecord)
}

func AddUsersLogs(db *gorm.DB, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	logRecord := []db_models.UsersLogs{}
	users := []db_models.Users{}
	err := FindRecords(db, &users, "users", filterFields, filterValues...)
	if err != nil {
		return err
	}
	logRecord, err = UsersToUsersLogs(users)
	fmt.Println(logRecord)
	if err != nil {
		return err
	}
	return AddRecords(db, logRecord)
}

func AddOrdersLogs(db *gorm.DB, userId int, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	logRecord := []db_models.OrdersLogs{}
	orders := []db_models.Orders{}
	err := FindRecords(db, &orders, "orders", filterFields, filterValues...)
	if err != nil {
		return err
	}
	logRecord, err = OrdersToOrdersLogs(userId, orders)
	if err != nil {
		return err
	}
	return AddRecords(db, logRecord)
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

func GetRoutines() []db_models.Routines {
	var ret []db_models.Routines
	database.DB.Raw("SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM tkoh_oms.routines LEFT JOIN tkoh_oms.locations C ON tkoh_oms.routines.start_location_id = C.location_id  LEFT JOIN tkoh_oms.locations D ON tkoh_oms.routines.end_location_id = D.location_id ").Scan(&ret)
	return ret
}
