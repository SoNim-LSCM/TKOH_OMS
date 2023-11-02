package orderManagement

import "tkoh_oms/models"

type UpdateDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
