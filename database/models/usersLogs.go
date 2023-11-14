package db_models

type UsersLogs struct {
	Id              int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	UserId          int    `json:"userId"`
	Username        string `json:"username" gorm:"type:string;column:username"`
	Password        string `json:"password" gorm:"type:string;column:password"`
	UserType        string `json:"userType" gorm:"type:string;column:user_type"`
	Token           string `json:"token"`
	TokenExpiryTime string `json:"tokenExpiryTime" gorm:"type:date;column:token_expiry_time"`
	LastLoginTime   string `json:"lastLoginTime" gorm:"type:date;column:last_login_time"`
	LastLogoutTime  string `json:"lastLogoutTime" gorm:"type:date;column:last_logout_time"`
	CreateTime      string `json:"createTime" gorm:"type:date;column:create_time"`
	LastUpdateTime  string `json:"lastUpdateTime" gorm:"type:date;column:last_update_time"`
	DutyLocationId  int    `json:"dutyLocationId"`
}
