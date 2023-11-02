package orderManagement

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type AddDeliveryOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
