package models

type FailResponseHeader struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	FailedReason    string `json:"failedReason"`
}
