package models

type FailResponse struct {
	Header FailResponseHeader `json:"header"`
}

func GetFailResponse(reason string) FailResponse {
	return FailResponse{Header: FailResponseHeader{ResponseCode: 400, ResponseMessage: "FAILED", FailedReason: reason}}
}
