package orderManagement

import "tkoh_oms/models"

type GetDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
