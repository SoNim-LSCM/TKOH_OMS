package ws_model

import "github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"

type WebsocketUpdateRoutineResponse struct {
	MessageCode string                           `json:"messageCode"`
	RoutineList orderManagement.RoutineOrderList `json:"routineList"`
}

func GetUpdateRoutineResponse(routineList orderManagement.RoutineOrderList) WebsocketUpdateRoutineResponse {
	return WebsocketUpdateRoutineResponse{MessageCode: "ORDER_UPDATE", RoutineList: routineList}
}
