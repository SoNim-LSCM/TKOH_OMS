package db_models

type Schedules struct {
	ScheduleID         int    `json:"scheduleId"`
	ScheduleStatus     string `json:"scheduleStatus"`
	ScheduleCraeteTime string `json:"scheduleCreateTime" gorm:"type:date;column:schedule_create_time"`
}
