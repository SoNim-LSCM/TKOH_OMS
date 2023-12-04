package dto

type ReportJobStatusDTO struct {
	JobID            int    `json:"jobId"`
	Status           string `json:"status"`
	Est              string `json:"est"`
	Eta              string `json:"eta"`
	ProcessingStatus string `json:"processingStatus"`
	Zone             string `json:"zone"`
	Location         string `json:"location"`
	RobotID          string `json:"robotId"`
	PayloadID        string `json:"payloadId"`
	MessageTime      string `json:"messageTime"`
}
