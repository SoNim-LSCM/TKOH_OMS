package orderManagement

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type AddRoutineResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   RoutineOrderListBody  `json:"body"`
}
