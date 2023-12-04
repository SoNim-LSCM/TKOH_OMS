package dto

type ReportJobStatusResponseDTO struct {
	ResponseCode    int                `json:"responseCode"`
	ResponseMessage string             `json:"responseMessage"`
	FailedReason    string             `json:"failedReason"`
	Body            ReportJobStatusDTO `json:"body"`
}
