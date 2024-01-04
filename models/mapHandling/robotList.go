package mapHandling

type RobotList []struct {
	RobotID           string    `json:"robotId"`
	RobotState        string    `json:"robotState"`
	RobotStatus       []string  `json:"robotStatus"`
	Zone              string    `json:"zone"`
	RobotPosition     []float64 `json:"robotPosition"`
	RobotOrientation  []float64 `json:"robotOrientation"`
	RobotCoordination []int     `json:"robotCoordination"`
	BatteryLevel      int       `json:"batteryLevel"`
	LastReportTime    string    `json:"lastReportTime"`
}
