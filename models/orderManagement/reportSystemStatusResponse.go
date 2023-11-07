package orderManagement

type ReportSystemStatusResponse struct {
	SystemState  string   `json:"systemState"`
	SystemStatus []string `json:"systemStatus"`
}
