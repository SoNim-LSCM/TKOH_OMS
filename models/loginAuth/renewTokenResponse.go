package loginAuth

import "github.com/SoNim-LSCM/TKOH_OMS/models"

type RenewTokenResponse struct {
	Header models.ResponseHeader  `json:"header"`
	Body   RenewTokenResponseBody `json:"body"`
}
