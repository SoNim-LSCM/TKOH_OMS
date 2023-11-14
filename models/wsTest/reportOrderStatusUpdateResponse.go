package wsTest

type ReportOrderStatusUpdateResponse struct {
	MessageCode      string   `json:"messageCode"`
	ScheduleID       int      `json:"scheduleId"`
	OrderID          int      `json:"orderId"`
	RobotID          []string `json:"robotId"`
	PayloadID        string   `json:"payloadId"`
	OrderStatus      string   `json:"orderStatus"`
	ProcessingStatus string   `json:"processingStatus"`
}
