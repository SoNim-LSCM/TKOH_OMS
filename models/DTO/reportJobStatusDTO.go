package dto

type ReportJobStatusDTO struct {
	JobId            int    `json:"jobId"`
	Status           string `json:"status"`
	Est              string `json:"est"`
	ETA              string `json:"eta"`
	ProcessingStatus string `json:"processingStatus"`
	Zone             string `json:"zone"`
	Location         string `json:"location"`
	RobotId          string `json:"robotId"`
	MessageTime      string `json:"messageTime"`
}
