package db_models

type OrdersLogs struct {
	Id                   int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ScheduleID           int    `json:"scheduleId"`
	OrderID              int    `json:"orderId"`
	OrderType            string `json:"orderType"`
	OrderCreatedType     string `json:"orderCreatedType"`
	OrderCreatedBy       int    `json:"orderCreatedBy"`
	OrderStatus          string `json:"orderStatus"`
	OrderStartTime       string `json:"startTime" gorm:"type:date;column:order_start_time"`
	ActualArrivalTime    string `json:"actualArrivalTime" gorm:"type:date;column:actual_arrival_time"`
	StartLocationID      int    `json:"startLocationId"`
	EndLocationID        int    `json:"endLocationId"`
	ExpectedStartTime    string `json:"expectedStartTime" gorm:"type:date;column:expected_start_time"`
	ExpectedDeliveryTime string `json:"expectedDeliveryTime" gorm:"type:date;column:expected_delivery_time"`
	ExpectedArrivalTime  string `json:"expectedArrivalTime" gorm:"type:date;column:expected_arrival_time"`
	ProcessingStatus     string `json:"processingStatus"`
	FailedReason         string `json:"failedReason"`
	LastUpdateTime       string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	LastUpdateBy         int    `json:"lastUpdateBy"`
}
