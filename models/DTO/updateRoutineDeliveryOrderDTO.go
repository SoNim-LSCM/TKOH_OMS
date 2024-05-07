package dto

type UpdateRoutineDeliveryOrderDTO struct {
	RoutineID      int    `json:"routineId"`
	OrderType      string `json:"orderType"`
	RoutinePattern struct {
		Day   []int `json:"day"`
		Month []int `json:"month"`
		Week  []int `json:"week"`
	} `json:"routinePattern"`
	NumberOfAmrRequire   int    `json:"numberOfAmrRequire"`
	StartLocationID      int    `json:"startLocationId"`
	StartLocationName    string `json:"startLocationName"`
	ExpectedStartTime    string `json:"expectedStartTime"`
	EndLocationID        int    `json:"endLocationId"`
	EndLocationName      string `json:"endLocationName"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime"`
	IsActive             bool   `json:"isActive"`
	PeriodStartTime      string `json:"periodStartTime"`
	PeriodEndTime        string `json:"periodEndTime"`
}
