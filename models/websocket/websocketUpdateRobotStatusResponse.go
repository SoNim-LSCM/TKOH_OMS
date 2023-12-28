package ws_model

import "github.com/SoNim-LSCM/TKOH_OMS/models/mapHandling"

type WebsocketUpdateRobotStatusResponse struct {
	MessageCode string                `json:"messageCode"`
	RobotList   mapHandling.RobotList `json:"robotList"`
}
