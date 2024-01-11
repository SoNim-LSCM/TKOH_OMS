package mapHandling

import (
	"math"

	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
)

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

func (robotList RobotList) CalculateCoordination(floors []db_models.Floors) (newList RobotList) {
	newList = robotList
	for i, robot := range robotList {
		for _, floor := range floors {
			if floor.FloorName == robot.Zone {
				quaternion := Quaternion{X: robot.RobotOrientation[0], Y: robot.RobotOrientation[1], Z: robot.RobotOrientation[2], W: robot.RobotOrientation[3]}
				phi, theta, psi := quaternion.Euler()
				euler := []float64{phi, theta, psi}
				newList[i].RobotOrientation = euler
				newList[i].RobotCoordination = GetCoordination(robot.RobotPosition, floor.OriginX, floor.OriginY, floor.Resolution, floor.MapX, floor.MapY)
				break
			}
		}
	}
	return
}

func GetCoordination(robotPosition []float64, originX float64, originY float64, resolution float64, mapX int, mapY int) []int {
	ret := []int{mapX + int((originX-robotPosition[0])/resolution), 0 - int((originY-robotPosition[1])/resolution)}
	return ret
}

// Quaternion represents a quaternion with a scalar and a vector part
type Quaternion struct {
	W, X, Y, Z float64
}

// FromEuler converts Euler angles (in radians) to a quaternion
func FromEuler(phi, theta, psi float64) Quaternion {
	c1 := math.Cos(phi / 2)
	s1 := math.Sin(phi / 2)
	c2 := math.Cos(theta / 2)
	s2 := math.Sin(theta / 2)
	c3 := math.Cos(psi / 2)
	s3 := math.Sin(psi / 2)
	return Quaternion{
		W: c1*c2*c3 + s1*s2*s3,
		X: s1*c2*c3 - c1*s2*s3,
		Y: c1*s2*c3 + s1*c2*s3,
		Z: c1*c2*s3 - s1*s2*c3,
	}
}

// Euler converts a quaternion to Euler angles (in radians)
func (q Quaternion) Euler() (phi, theta, psi float64) {
	phi = math.Atan2(2*(q.W*q.X+q.Y*q.Z), 1-2*(q.X*q.X+q.Y*q.Y))
	theta = math.Asin(2 * (q.W*q.Y - q.Z*q.X))
	psi = math.Atan2(2*(q.W*q.Z+q.X*q.Y), 1-2*(q.Y*q.Y+q.Z*q.Z))
	return
}
