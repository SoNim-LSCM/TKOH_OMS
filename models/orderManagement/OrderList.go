package orderManagement

type OrderList []struct {
	ScheduleID            int    `json:"scheduleId"`
	OrderID               int    `json:"orderId"`
	OrderType             string `json:"orderType"`
	OrderCreatedType      string `json:"orderCreatedType"`
	OrderCreatedBy        int    `json:"orderCreatedBy"`
	OrderStatus           string `json:"orderStatus"`
	StartTime             string `json:"startTime"`
	EndTime               string `json:"endTime"`
	StartLocationID       int    `json:"startLocationId"`
	StartLocationName     string `json:"startLocationName"`
	ExpectingStartTime    string `json:"expectingStartTime"`
	EndLocationID         int    `json:"endLocationId"`
	EndLocationName       string `json:"endLocationName"`
	ExpectingDeliveryTime string `json:"expectingDeliveryTime"`
	ProcessingStatus      string `json:"processingStatus"`
}
