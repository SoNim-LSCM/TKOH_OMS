package db_models

type Orders struct {
	ScheduleID         int    `json:"scheduleId"`
	OrderID            int    `json:"orderId"`
	OrderType          string `json:"orderType"`
	OrderCreatedType   string `json:"orderCreatedType"`
	OrderCreatedBy     int    `json:"orderCreatedBy"`
	OrderStatus        string `json:"orderStatus"`
	OrderStartTime     string `json:"startTime" gorm:"type:date;column:order_start_time"`
	OrderEndTime       string `json:"endTime" gorm:"type:date;column:order_end_time"`
	ActualArrivalTime  string `json:"actualArrivalTime" gorm:"type:date;column:actual_arrival_time"`
	StartLocationID    int    `json:"startLocationId"`
	EndLocationID      int    `json:"endLocationId"`
	ExpectStartTime    string `json:"expectingStartTime" gorm:"type:date;column:expect_start_time"`
	ExpectDeliveryTime string `json:"expectingDeliveryTime" gorm:"type:date;column:expect_delivey_time"`
	ExpectArrivalTime  string `json:"expectingArrivalTime" gorm:"type:date;column:expect_arrival_time"`
	ProcessingStatus   string `json:"processingStatus"`
	FailedReason       string `json:"failedReason"`
}
