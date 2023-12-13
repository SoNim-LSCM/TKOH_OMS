package db_models

type JobsLogs struct {
	ID                  int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	CreateID            int    `json:"createId"`
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
}
