package dto

import "tkoh_oms/models/mapHandling"

type UpdateRobotStatusDTOResponse struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Body            struct {
		RobotList mapHandling.RobotList `json:"robotList"`
	} `json:"body"`
}
