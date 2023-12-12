package dto

type ReportJobStatusResponseDTO struct {
	ResponseCode    int                `json:"responseCode"`
	ResponseMessage string             `json:"responseMessage"`
	FailReason      string             `json:"failReason"`
	Body            ReportJobStatusDTO `json:"body"`
}
