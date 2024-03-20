package loginAuth

import "tkoh_oms/models"

type LoginResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   LoginResponseBody     `json:"body"`
}
