package orderManagement

import "tkoh_oms/models"

type UpdateRoutineDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   RoutineOrderListBody  `json:"body"`
}
