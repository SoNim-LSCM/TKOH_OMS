package rfms

import dto "tkoh_oms/models/DTO"

type CreateJobResponse struct {
	ResponseCode    int                    `json:"responseCode"`
	ResponseMessage string                 `json:"responseMessage"`
	Body            dto.ReportJobStatusDTO `json:"body"`
}
