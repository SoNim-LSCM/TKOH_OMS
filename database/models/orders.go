package db_models

type Orders struct {
	ScheduleID           int    `json:"scheduleId"`
	OrderID              int    `json:"orderId" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	OrderType            string `json:"orderType"`
	OrderCreatedType     string `json:"orderCreatedType"`
	OrderCreatedBy       int    `json:"orderCreatedBy"`
	OrderStatus          string `json:"orderStatus"`
	OrderStartTime       string `json:"startTime" gorm:"type:date;column:order_start_time"`
	ActualArrivalTime    string `json:"actualArrivalTime" gorm:"type:date;column:actual_arrival_time"`
	StartLocationID      int    `json:"startLocationId"`
	EndLocationID        int    `json:"endLocationId"`
	ExpectedStartTime    string `json:"expectedStartTime" gorm:"type:date;column:expected_start_time"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime" gorm:"type:date;column:expected_delivey_time"`
	ExpectedArrivalTime  string `json:"expectedArrivalTime" gorm:"type:date;column:expected_arrival_time"`
	ProcessingStatus     string `json:"processingStatus"`
	FailedReason         string `json:"failedReason"`
}
