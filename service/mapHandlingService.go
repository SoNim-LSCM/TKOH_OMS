package service

import (
	"encoding/json"
	"errors"
	"log"

	apiHandler "github.com/SoNim-LSCM/TKOH_OMS/api"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/mapHandling"
	"github.com/SoNim-LSCM/TKOH_OMS/models/rfms"
	ws_model "github.com/SoNim-LSCM/TKOH_OMS/models/websocket"
	"github.com/SoNim-LSCM/TKOH_OMS/websocket"
	"gorm.io/gorm"
)

func FindAllDutyRooms() ([]db_models.Locations, error) {
	var val []db_models.Locations
	if database.CheckDatabaseConnection() {
		err := database.DB.Find(&val).Error
		return val, err
	} else {
		return val, errors.New("Database Connection Fail")
	}
	return val, nil
}

func GetFloorPlan() ([]db_models.Floors, error) {
	database.CheckDatabaseConnection()
	var val []db_models.Floors
	err := database.DB.Find(&val).Error
	return val, err
}

func FloorPlanToMapList(floorPlan []db_models.Floors) (mapHandling.MapList, error) {
	mapList := mapHandling.MapList{}
	bJson, err := json.Marshal(floorPlan)
	if err != nil {
		return mapList, err
	}
	err = json.Unmarshal(bJson, &mapList)
	if err != nil {
		return mapList, err
	}
	return mapList, err
}

func GetLocationFromRFMS() error {

	response, err := apiHandler.Get("/locationList?type=DESTINATION", nil)
	if err != nil {
		return err
	}
	jsonResponse := rfms.GetLocationResponse{}
	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		return err
	}
	locations, err := locationListToDBLocations(jsonResponse.Body.LocationList)
	if err != nil {
		return err
	}

	log.Print(locations)

	if !database.CheckDatabaseConnection() {
		return errors.New("Cannot Connect DB")
	}
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err := TruncateTable(tx, "locations")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err := AddRecords(tx, locations)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func BackgroundReportRobotStatus() error {

	response, err := apiHandler.Get("/robotStatus?robotType=AMR", nil)
	if err != nil {
		return err
	}

	updateJobStatus := dto.UpdateRobotStatusDTOResponse{}
	err = json.Unmarshal(response, &updateJobStatus)
	if err != nil {
		return err
	}
	if updateJobStatus.ResponseCode != 200 {
		return errors.New("Get Robot Status from RFMS Failed")
	}

	log.Printf("BackgroundReportRobotStatus: %s", ws_model.GetUpdateRobotResponse(updateJobStatus.Body.RobotList))

	websocket.SendBoardcastMessage(ws_model.GetUpdateRobotResponse(updateJobStatus.Body.RobotList))

	return nil
}

func locationListToDBLocations(locationList interface{}) ([]db_models.Locations, error) {
	locations := []db_models.Locations{}
	bJson, err := json.Marshal(locationList)
	if err != nil {
		return locations, err
	}
	err = json.Unmarshal(bJson, &locations)
	if err != nil {
		return locations, err
	}
	return locations, err
}
