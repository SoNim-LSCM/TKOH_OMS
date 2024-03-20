package orderManagement

import "tkoh_oms/models"

type GetRoutineDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   RoutineOrderListBody  `json:"body"`
}
