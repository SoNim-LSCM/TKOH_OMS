package db_models

type SchedulesLogs struct {
	Id                 int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ScheduleID         int    `json:"scheduleId"`
	ScheduleStatus     string `json:"scheduleStatus"`
	ScheduleCraeteTime string `json:"scheduleCreateTime" gorm:"type:date;column:schedule_create_time"`
	OrderType          string `json:"orderType"`
	OrderCreatedType   string `json:"orderCreatedType"`
	NumberOfAmrRequire int    `json:"numberOfAmrRequire"`
	RoutineID          int    `json:"routineId"`
	LastUpdateTime     string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	LastUpdateBy       int    `json:"lastUpdateBy"`
}
