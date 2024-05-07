package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"tkoh_oms/database"
	db_models "tkoh_oms/database/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TruncateTable(db *gorm.DB, tableName string) error {
	database.CheckDatabaseConnection()
	// query := "SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM "+os.Getenv("MYSQL_DB_NAME")+"." + table + " LEFT JOIN "+os.Getenv("MYSQL_DB_NAME")+".locations C ON "+os.Getenv("MYSQL_DB_NAME")+"." + table + ".start_location_id = C.location_id  LEFT JOIN "+os.Getenv("MYSQL_DB_NAME")+".locations D ON "+os.Getenv("MYSQL_DB_NAME")+"." + table + ".end_location_id = D.location_id WHERE " + filterFields
	err := db.Exec("truncate table " + tableName).Error
	if err != nil {
		return err
	}
	return nil
}

func FindRecordsWithRaw(db *gorm.DB, records interface{}, query string, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	// query := "SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM "+os.Getenv("MYSQL_DB_NAME")+"." + table + " LEFT JOIN "+os.Getenv("MYSQL_DB_NAME")+".locations C ON "+os.Getenv("MYSQL_DB_NAME")+"." + table + ".start_location_id = C.location_id  LEFT JOIN "+os.Getenv("MYSQL_DB_NAME")+".locations D ON "+os.Getenv("MYSQL_DB_NAME")+"." + table + ".end_location_id = D.location_id WHERE " + filterFields
	err := db.Raw(query, filterValues...).Scan(records).Error
	if err != nil {
		return err
	}
	return nil
}

func FindRecords(db *gorm.DB, records interface{}, table string, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	if err := db.Table(table).Where(filterFields, filterValues...).Find(records).Error; err != nil {
		return err
	}
	return nil
}

func UpdateRecords(db *gorm.DB, updatedRecordList interface{}, table string, updateMap map[string]interface{}, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Table(table).Where(filterFields, filterValues...).Updates(updateMap)
	err := FindRecords(db, updatedRecordList, table, filterFields, filterValues...)
	if err != nil {
		return err
	}
	return nil
}

func AddRecords(db *gorm.DB, records interface{}) error {
	database.CheckDatabaseConnection()
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

func AddRoutinesLogs(db *gorm.DB, userId int, filterFields interface{}, filterValues ...interface{}) error {
	database.CheckDatabaseConnection()
	logRecord := []db_models.RoutinesLogs{}
	routines := []db_models.Routines{}
	err := FindRecords(db, &routines, "routines", filterFields, filterValues...)
	if err != nil {
		return err
	}
	logRecord, err = RoutinesToRoutinesLogs(userId, routines)
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

func GetRoutines() []db_models.Routines {
	var ret []db_models.Routines
	database.DB.Raw("SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM " + os.Getenv("MYSQL_DB_NAME") + ".routines LEFT JOIN " + os.Getenv("MYSQL_DB_NAME") + ".locations C ON " + os.Getenv("MYSQL_DB_NAME") + ".routines.start_location_id = C.location_id  LEFT JOIN " + os.Getenv("MYSQL_DB_NAME") + ".locations D ON " + os.Getenv("MYSQL_DB_NAME") + ".routines.end_location_id = D.location_id ").Scan(&ret)
	return ret
}
