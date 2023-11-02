package systemStatus

type SystemStatusResponse struct {
	MessageCode  string   `json:"messageCode"`
	SystemState  string   `json:"systemState"`
	SystemStatus []string `json:"systemStatus"`
}
