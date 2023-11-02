package db_models

import "time"

type Users struct {
	UserId              string    `json:"user_id"`
	Username            string    `json:"username"`
	Password            string    `json:"password"`
	UserType            string    `json:"user_type"`
	Token               string    `json:"token"`
	TokenExpiryDateTime time.Time `json:"token_expiry_datetime"`
	LastLoginDateTime   time.Time `json:"last_login_datetime"`
	LastLogoutDateTime  time.Time `json:"last_logout_datetime"`
	CreateDateTime      time.Time `json:"create_datetime"`
	LastUpdateDateTime  time.Time `json:"lastUpdate_datetime"`
}
