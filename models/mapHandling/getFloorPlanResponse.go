package mapHandling

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type GetFloorPlanResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   MapListBody           `json:"body"`
}
