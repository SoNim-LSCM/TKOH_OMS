package orderManagement

type RoutineOrderList []struct {
	RoutineID            int            `json:"routineId"`
	OrderType            string         `json:"orderType"`
	RoutinePattern       RoutinePattern `json:"routinePattern"`
	NextDeliveryDate     string         `json:"nextDeliveryDate"`
	IsActive             bool           `json:"isActive"`
	RoutineCreatedBy     int            `json:"routineCreatedBy"`
	NumberOfAmrRequire   int            `json:"numberOfAmrRequire"`
	StartLocationID      int            `json:"startLocationId"`
	StartLocationName    string         `json:"startLocationName"`
	ExpectedStartTime    string         `json:"expectedStartTime"`
	EndLocationID        int            `json:"endLocationId"`
	EndLocationName      string         `json:"endLocationName"`
	ExpectedDeliveryTime string         `json:"expectedDeliveryTime"`
}
