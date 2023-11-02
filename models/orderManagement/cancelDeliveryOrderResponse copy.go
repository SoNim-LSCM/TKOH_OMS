package orderManagement

import "tkoh_oms/models"

type CancelDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
