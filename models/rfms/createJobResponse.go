package rfms

import dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"

type CreateJobResponse struct {
	ResponseCode    int                    `json:"responseCode"`
	ResponseMessage string                 `json:"responseMessage"`
	Body            dto.ReportJobStatusDTO `json:"body"`
}
