package wsTest

type ReportDutyLocationUpdateResponse struct {
	MessageCode      string `json:"messageCode"`
	UserID           int    `json:"userId"`
	DutyLocationID   int    `json:"dutyLocationId"`
	DutyLocationName string `json:"dutyLocationName"`
}
