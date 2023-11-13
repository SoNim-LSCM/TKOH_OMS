package failReason

type FAIL_REASON string

const (
	SystemNotReady FAIL_REASON = "SYSTEM_NOT_READY"
	RobotNotReady  FAIL_REASON = "ROBOT_NOT_READY"
)
