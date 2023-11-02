package orderManagement

import "tkoh_oms/models"

type AddDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
