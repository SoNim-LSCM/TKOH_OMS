package orderManagement

import "tkoh_oms/models"

type TriggerHandlingOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
