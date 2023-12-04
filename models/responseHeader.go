package models

type ResponseHeader struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	FailedReason    string `json:"failedReason"`
}

func GetSuccessResponseHeader() ResponseHeader {
	return ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
}
