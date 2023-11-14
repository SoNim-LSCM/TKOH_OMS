package db_models

type RoutinesLogs struct {
	ID                   int    `json:"Id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	RoutineID            int    `json:"routineId"`
	OrderType            string `json:"orderType"`
	RoutinePattern       string `json:"routinePattern"`
	IsActive             bool   `json:"isActive" gorm:"type:boolean"`
	NumberOfAmrRequire   int    `json:"numberOfAmrRequire"`
	StartLocationID      int    `json:"startLocationId"`
	ExpectedStartTime    string `json:"expectedStartTime"`
	EndLocationID        int    `json:"endLocationId"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime"`
	LastUpdateTime       string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	LastUpdateBy         int    `json:"lastUpdateBy"`
}
