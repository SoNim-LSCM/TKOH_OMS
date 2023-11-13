package orderManagement

type OrderList []struct {
	ScheduleID           int    `json:"scheduleId"`
	OrderID              int    `json:"orderId"`
	OrderType            string `json:"orderType"`
	OrderCreatedType     string `json:"orderCreatedType"`
	OrderCreatedBy       int    `json:"orderCreatedBy"`
	OrderStatus          string `json:"orderStatus"`
	StartTime            string `json:"startTime"`
	EndTime              string `json:"endTime"`
	ActualArrivalTime    string `json:"actualArrivalTime"`
	StartLocationID      int    `json:"startLocationId"`
	StartLocationName    string `json:"startLocationName"`
	ExpectedStartTime    string `json:"expectedStartTime"`
	EndLocationID        int    `json:"endLocationId"`
	EndLocationName      string `json:"endLocationName"`
	ExpectedArrivalTime  string `json:"expectedArrivalTime"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime"`
	ProcessingStatus     string `json:"processingStatus"`
	FailedReason         string `json:"failedReason"`
}
