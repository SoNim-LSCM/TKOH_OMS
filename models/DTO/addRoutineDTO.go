package dto

type AddRoutineDTO struct {
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
	PeriodStartTime      string `json:"periodStartTime"`
	PeriodEndTime        string `json:"periodEndTime"`
}
