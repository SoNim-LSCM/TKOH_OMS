package db_models

type Schedules struct {
	ScheduleID         int    `json:"scheduleId" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ScheduleStatus     string `json:"scheduleStatus"`
	ScheduleCraeteTime string `json:"scheduleCreateTime" gorm:"type:date;column:schedule_create_time"`
	OrderType          string `json:"orderType"`
	NumberOfAmrRequire int    `json:"numberOfAmrRequire"`
	LastUpdateTime     string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	LastUpdateBy       int    `json:"lastUpdateBy"`
}
