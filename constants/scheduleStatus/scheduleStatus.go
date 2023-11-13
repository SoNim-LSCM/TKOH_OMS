package scheduleStatus

type SCHEDULE_STATUS string

const (
	Waiting           SCHEDULE_STATUS = "WAITING"
	ReadyForLoading   SCHEDULE_STATUS = "READY_FOR_LOADING"
	ReadyForUnloading SCHEDULE_STATUS = "READY_FOR_UNLOADING"
	Incoming          SCHEDULE_STATUS = "INCOMING"
	Error             SCHEDULE_STATUS = "ERROR"
	Loading           SCHEDULE_STATUS = "LOADING"
	Unloading         SCHEDULE_STATUS = "UNLOADING"
	Completed         SCHEDULE_STATUS = "COMPLETED"
)
