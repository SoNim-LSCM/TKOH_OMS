package rfms

type StatusResponse struct {
	JobID            int    `json:"jobId"`
	Status           string `json:"status"`
	Est              int64  `json:"est"`
	Eta              int64  `json:"eta"`
	ProcessingStatus string `json:"processingStatus"`
	Zone             string `json:"zone"`
	Location         string `json:"location"`
	RobotID          string `json:"robotId"`
	PayloadID        string `json:"payloadId"`
	MessageTime      int64  `json:"messageTime"`
}
