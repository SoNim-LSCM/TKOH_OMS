package rfms

type CreateJobRequest struct {
	JobNature  string `json:"jobNature"`
	LocationID int    `json:"locationId"`
	RobotID    string `json:"robotId"`
	PayloadID  string `json:"payloadId"`
}
