package loginAuth

import "tkoh_oms/models"

type RenewTokenResponse struct {
	Header models.ResponseHeader  `json:"header"`
	Body   RenewTokenResponseBody `json:"body"`
}
