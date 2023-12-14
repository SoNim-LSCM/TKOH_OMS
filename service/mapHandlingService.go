package service

import (
	"encoding/json"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	"github.com/SoNim-LSCM/TKOH_OMS/models/mapHandling"
)

func FindAllDutyRooms() ([]db_models.Locations, error) {
	database.CheckDatabaseConnection()
	var val []db_models.Locations
	err := database.DB.Find(&val).Error
	return val, err
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
