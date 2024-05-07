package dto

type AddDeliveryOrderDTO struct {
	OrderType            string `json:"orderType"`
	NumberOfAmrRequire   int    `json:"numberOfAmrRequire"`
	StartLocationID      int    `json:"startLocationId"`
	StartLocationName    string `json:"startLocationName"`
	ExpectedStartTime    string `json:"expectedStartTime"`
	RoutineID            int    `json:"routineId"`
	EndLocationID        int    `json:"endLocationId"`
	EndLocationName      string `json:"endLocationName"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime"`
	RobotId              string `json:"robotId"`
}
