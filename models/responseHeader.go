package models

type ResponseHeader struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

func GetSuccessResponseHeader() ResponseHeader {
	return ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
}
