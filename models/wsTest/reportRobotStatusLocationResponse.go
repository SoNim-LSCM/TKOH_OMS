package wsTest

type ReportRobotStatusLocationResponse struct {
	MessageCode string `json:"messageCode"`
	RobotList   []struct {
		RobotID             string    `json:"robotId"`
		RobotCoordatination []int     `json:"robotCoordatination"`
		RobotPostion        []float64 `json:"robotPostion"`
		RobotOrienation     []float64 `json:"robotOrienation"`
		RobotState          string    `json:"robotState"`
		RobotStatus         []string  `json:"robotStatus"`
		BatteryLevel        float64   `json:"batteryLevel"`
		LastReportTime      string    `json:"lastReportTime"`
	} `json:"robotList"`
}
