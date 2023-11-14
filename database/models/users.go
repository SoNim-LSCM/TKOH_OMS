package db_models

type Users struct {
	UserId          int    `json:"userId" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Username        string `json:"username" gorm:"type:string;column:username"`
	Password        string `json:"password" gorm:"type:string;column:password"`
	UserType        string `json:"userType" gorm:"type:string;column:user_type"`
	Token           string `json:"token"`
	TokenExpiryTime string `json:"tokenExpiryTime" gorm:"type:date;column:token_expiry_time"`
	LastLoginTime   string `json:"lastLoginTime" gorm:"type:date;column:last_login_time"`
	LastLogoutTime  string `json:"lastLogoutTime" gorm:"type:date;column:last_logout_time"`
	CreateTime      string `json:"createTime" gorm:"type:date;column:create_time"`
	DutyLocationId  int    `json:"dutyLocationId"`
	LastUpdateTime  string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
}
