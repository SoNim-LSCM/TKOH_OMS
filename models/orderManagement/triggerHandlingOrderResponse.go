package orderManagement

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type TriggerHandlingOrderResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   OrderListBody         `json:"body"`
}
