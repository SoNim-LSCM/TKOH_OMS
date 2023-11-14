package db_models

type Routines struct {
	RoutineID            int    `json:"routineId" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	OrderType            string `json:"orderType"`
	RoutinePattern       string `json:"routinePattern"`
	IsActive             bool   `json:"isActive" gorm:"type:boolean"`
	NumberOfAmrRequire   int    `json:"numberOfAmrRequire"`
	StartLocationID      int    `json:"startLocationId"`
	StartLocationName    string `json:"startLocationName" gorm:"<-:false"`
	EndLocationID        int    `json:"endLocationId"`
	EndLocationName      string `json:"endLocationName" gorm:"<-:false"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime" gorm:"type:date;column:expected_delivey_time"`
	LastUpdateTime       string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	LastUpdateBy         int    `json:"lastUpdateBy"`
}
