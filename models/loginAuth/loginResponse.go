package loginAuth

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type LoginResponse struct {
	Header models.ResponseHeader `json:"header"`
	Body   LoginResponseBody     `json:"body"`
}
