package mapHandling

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type GetDutyRoomsResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   LocationListBody      `json:"body"`
}
