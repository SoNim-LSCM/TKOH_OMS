package service

import (
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
)

func FindAllDutyRooms() ([]db_models.Locations, error) {
	database.CheckDatabaseConnection()
	var val []db_models.Locations
	err := database.DB.Find(&val).Error
	return val, err
}
