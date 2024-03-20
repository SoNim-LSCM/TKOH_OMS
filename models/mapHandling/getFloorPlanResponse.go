package mapHandling

import "tkoh_oms/models"

type GetFloorPlanResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   MapListBody           `json:"body"`
}
