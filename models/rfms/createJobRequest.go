package rfms

type CreateJobRequest struct {
	JobNature string `json:"jobNature"`
	Zone      string `json:"zone"`
	Location  string `json:"location"`
}
