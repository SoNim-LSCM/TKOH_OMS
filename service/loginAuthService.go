package service

import (
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
)

func UpdateUser(username string, userType string, updateFields []string, updateValues ...interface{}) ([]db_models.Users, error) {
	updatedUser := []db_models.Users{}
	updateMap := utils.CreateMap(updateFields, updateValues...)
	err := UpdateRecords(&updatedUser, "users", updateMap, "username = ? AND user_type = ?", username, userType)
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
