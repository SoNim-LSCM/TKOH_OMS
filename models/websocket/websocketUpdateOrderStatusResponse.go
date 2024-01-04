package ws_model

type WebsocketUpdateOrderStatusResponse struct {
	MessageCode      string   `json:"messageCode"`
	OrderID          int      `json:"orderId"`
	OrderStatus      string   `json:"orderStatus"`
	PayloadID        string   `json:"payloadId"`
	ProcessingStatus string   `json:"processingStatus"`
	RobotID          []string `json:"robotId"`
	ScheduleID       int      `json:"scheduleId"`
}

func GetUpdateOrderStatusResponse(orderId int, orderStatus string, payloadID string, processingStatus string, robotID []string, scheduleID int) WebsocketUpdateOrderStatusResponse {
	return WebsocketUpdateOrderStatusResponse{MessageCode: "ORDER_STATUS", OrderID: orderId, OrderStatus: orderStatus, PayloadID: payloadID, ProcessingStatus: processingStatus, RobotID: robotID, ScheduleID: scheduleID}
}
