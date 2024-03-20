package mapHandling

import "tkoh_oms/models"

type GetDutyRoomsResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   LocationListBody      `json:"body"`
}
