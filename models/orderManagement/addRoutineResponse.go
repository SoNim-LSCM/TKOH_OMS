package orderManagement

import "tkoh_oms/models"

type AddRoutineResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   RoutineOrderListBody  `json:"body"`
}
