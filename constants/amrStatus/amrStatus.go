package amrStatus

type AMR_STATUS string

const (
	Idle       AMR_STATUS = "IDLE"
	Charging   AMR_STATUS = "CHARGING"
	Delivering AMR_STATUS = "DELIVERING"
	Error      AMR_STATUS = "ERROR"
	Loading    AMR_STATUS = "LOADING"
	Unloading  AMR_STATUS = "UNLOADING"
)
