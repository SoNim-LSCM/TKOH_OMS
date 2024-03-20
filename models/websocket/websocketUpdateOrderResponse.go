package ws_model

import "tkoh_oms/models/orderManagement"

type WebsocketUpdateOrderResponse struct {
	MessageCode string                    `json:"messageCode"`
	OrderList   orderManagement.OrderList `json:"orderList"`
}

func GetUpdateOrderResponse(orderList orderManagement.OrderList) WebsocketUpdateOrderResponse {
	return WebsocketUpdateOrderResponse{MessageCode: "ORDER_UPDATE", OrderList: orderList}
}
