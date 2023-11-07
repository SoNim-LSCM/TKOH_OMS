package db_models

type Users struct {
	UserId              int    `json:"user_id"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	UserType            string `json:"user_type"`
	Token               string `json:"token"`
	TokenExpiryDateTime string `json:"token_expiry_datetime" gorm:"type:date;column:token_expiry_datetime"`
	LastLoginDateTime   string `json:"last_login_datetime" gorm:"type:date;column:last_login_datetime"`
	LastLogoutDateTime  string `json:"last_logout_datetime" gorm:"type:date;column:last_logout_datetime"`
	CreateDateTime      string `json:"create_datetime" gorm:"type:date;column:create_datetime"`
	LastUpdateDateTime  string `json:"lastUpdate_datetime" gorm:"type:date;column:last_update_datetime"`
	DutyLocationId      int    `json:"duty_location_id"`
}
