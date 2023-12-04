package models

type FailResponse struct {
	Header FailResponseHeader `json:"header"`
}

func GetFailResponse(description string, err string) FailResponse {
	return FailResponse{Header: FailResponseHeader{ResponseCode: 400, ResponseMessage: "FAILED", FailedReason: description + ": " + err}}
}
