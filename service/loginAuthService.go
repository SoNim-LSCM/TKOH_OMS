package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"

	"gorm.io/gorm/clause"
)

const (
	LOGIN  int = 0
	LOGOUT     = 1
	RENEW      = 2
)

func FindUser(username string, userType string) (db_models.Users, error) {
	database.CheckDatabaseConnection()
	user := db_models.Users{}
	val := make(map[string]interface{})
	log.Printf("mysql query: FindUser: %s, %s\n", username, userType)
	if err := database.DB.Table("users").Find(&val, "username = ? AND user_type = ?", username, userType).Error; err != nil {
		log.Printf("mysql query error: %s\n", err.Error())
	}
	if len(val) == 0 {
		return user, errors.New("Username not found")
	}
	jsonString, _ := json.Marshal(val)
	json.Unmarshal(jsonString, &user)
	return user, nil
}

func UpdateUserToken(username string, userType string, token string, tokenExpire int64, actionType int) (db_models.Users, error) {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: UpdateUserToken: %s, %s\n", username, userType)
	var updatedUser db_models.Users
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	// map[string]interface{}{"token_expiry_datetime": timeExpire, "lastLogin_datetime": timeNow, "token": token}

	timeExpire := time.Unix(tokenExpire, 0).Format("2006-01-02 15:04:05")
	switch actionType {
	case LOGIN:
		database.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Table("users").Where("username = ? AND user_type = ?", username, userType).Updates(map[string]interface{}{"token_expiry_datetime": timeExpire, "last_login_datetime": timeNow, "token": token})
	case LOGOUT:
		database.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Table("users").Where("username = ? AND user_type = ?", username, userType).Updates(map[string]interface{}{"last_logout_datetime": timeNow, "token": nil})
	case RENEW:
		database.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Table("users").Where("username = ? AND user_type = ?", username, userType).Updates(map[string]interface{}{"token_expiry_datetime": timeExpire, "token": token})
	default:
		return updatedUser, errors.New("Unexpected action")
	}

	updatedUser, err := FindUser(username, userType)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func UsersToLoginResponse(user db_models.Users) (loginAuth.LoginResponseBody, error) {
	lastLogin, err := StringToResponseTime(user.LastLoginDateTime)
	if err != nil {
		return loginAuth.LoginResponseBody{}, err
	}
	expireTime, err := StringToResponseTime(user.TokenExpiryDateTime)
	if err != nil {
		return loginAuth.LoginResponseBody{}, err
	}
	return loginAuth.LoginResponseBody{ID: user.UserId, Username: user.Username, UserType: user.UserType, AuthToken: user.Token, LoginDateTime: lastLogin, TokenExpiryDateTime: expireTime, DutyLocationId: user.DutyLocationId}, nil
}
