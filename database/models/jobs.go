package db_models

type Jobs struct {
	CreateID            int    `json:"createId" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	OrderID             int    `json:"orderId"`
	JobID               int    `json:"jobId" `
	JobType             string `json:"jobType"`
	JobStatus           string `json:"jobStatus"`
	ProcessingStatus    string `json:"processingStatus"`
	JobStartTime        string `json:"jobStartTime" gorm:"type:date;column:job_start_time"`
	ExpectedArrivalTime string `json:"expectedArrivalTime" gorm:"type:date;column:expected_arrival_time"`
	EndLocationID       int    `json:"endLocationId"`
	FailedReason        string `json:"failedReason"`
	LastUpdateTime      string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	StatusLocation      string `json:"status_location"`
	RobotID             string `json:"robotId"`
	PayloadID           string `json:"payloadId"`
}
