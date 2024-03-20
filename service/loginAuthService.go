package service

import (
	"encoding/json"
	"errors"
	"strings"

	"tkoh_oms/database"
	db_models "tkoh_oms/database/models"
	dto "tkoh_oms/models/DTO"
	"tkoh_oms/models/loginAuth"
	"tkoh_oms/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func FindUsers(filterFields interface{}, filterValues ...interface{}) ([]db_models.Users, error) {
	var users []db_models.Users
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := FindRecords(tx, &users, "users", filterFields, filterValues...); err != nil {
			return errors.New("Failed to search: " + err.Error())
		}
		return nil
	})
	return users, err
}

func UpdateUser(db *gorm.DB, username string, userType string, updateFields []string, updateValues ...interface{}) ([]db_models.Users, error) {
	updatedUser := []db_models.Users{}
	updateMap := utils.CreateMap(updateFields, updateValues...)
	err := AddUsersLogs(db, "username = ? AND user_type = ?", username, userType)
	if err != nil {
		return updatedUser, errors.New("Failed to create log")
		// return updatedUser, err
	}
	err = UpdateRecords(db, &updatedUser, "users", updateMap, "username = ? AND user_type = ?", username, userType)
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, err
}

func LoginStaff(request *dto.LoginStaffDTO) ([]db_models.Users, error) {
	var users []db_models.Users
	var updatedUser []db_models.Users
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := FindRecords(tx, &users, "users", &db_models.Users{Username: request.Username, UserType: "STAFF"})
		if err != nil {
			return err
		}

		if len(users) == 0 {
			return errors.New("User not found")
		}

		// if request.DutyLocationId != users[0].DutyLocationId {
		// 	return errors.New("Failed to search for staff")
		// }

		if isValid, err := utils.ValidateJwtToken(users[0].Token); err != nil || isValid {
			if isValid {
				return errors.New("this account is already logged in by someone")
			} else {
				if err != nil {
					if !strings.Contains(err.Error(), "token is expired") {
						return errors.New("failed to validate token: " + err.Error())
					}
				}
			}
		}

		token, expiryTime, err := utils.GenerateJwtStaff(users[0].UserId, users[0].Username, request.DutyLocationId)
		if err != nil {
			return err
		}

		updatedUser, err = UpdateUser(tx, users[0].Username, users[0].UserType, []string{"token_expiry_time", "last_login_time", "token", "duty_location_id"}, utils.TimeInt64ToString(expiryTime), utils.GetTimeNowString(), token, request.DutyLocationId)

		if err != nil {
			return err
		}

		return nil
	})
	return updatedUser, err
}

func LoginAdmin(request *dto.LoginAdminDTO) ([]db_models.Users, error) {
	var users []db_models.Users
	var updatedUser []db_models.Users
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := FindRecords(tx, &users, "users", &db_models.Users{Username: request.Username, UserType: "ADMIN"})
		if err != nil {
			return err
		}

		if len(users) == 0 {
			return errors.New("User not found")
		}

		if isValid, err := utils.ValidateJwtToken(users[0].Token); err != nil || isValid {
			if isValid {
				return errors.New("this account is already logged in by someone")
			} else {
				if err != nil {
					if !strings.Contains(err.Error(), "token is expired") {
						return errors.New("failed to validate token: " + err.Error())
					}
				}
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(request.Password))
		if err != nil {
			return err
		}

		token, expiryTime, err := utils.GenerateJwtAdmin(users[0].UserId, users[0].Username, users[0].Password)
		if err != nil {
			return err
		}

		updatedUser, err = UpdateUser(tx, users[0].Username, users[0].UserType, []string{"token_expiry_time", "last_login_time", "token"}, utils.TimeInt64ToString(expiryTime), utils.GetTimeNowString(), token)
		if err != nil {
			return err
		}

		return nil
	})
	return updatedUser, err
}

func Logout(claim *utils.Claims, token string) ([]db_models.Users, error) {

	var users []db_models.Users
	var updatedUser []db_models.Users
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := FindRecords(tx, &users, "users", &db_models.Users{Username: claim.Username, UserType: claim.UserType})
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return errors.New("User not found")
		}
		if users[0].Token == "" {
			return errors.New("Account logged out already")
		} else if users[0].Token != token {
			return errors.New("Incorrect token")
		}

		_, err = UpdateUser(tx, claim.Username, claim.UserType, []string{"last_logout_time", "token"}, utils.GetTimeNowString(), "")
		if err != nil {
			return err
		}

		return nil
	})
	return updatedUser, err
}

func RenewToken(claim *utils.Claims, token string) ([]db_models.Users, error) {

	var users []db_models.Users
	var updatedUser []db_models.Users
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := FindRecords(tx, &users, "users", &db_models.Users{Username: claim.Username, UserType: claim.UserType})
		if err != nil {
			return err
		}

		if len(users) == 0 {
			return errors.New("User not found")
		}
		if users[0].Token == "" {
			return errors.New("Account logged out already")
		} else if users[0].Token != token {
			return errors.New("Incorrect token")
		}

		if isValid, err := utils.ValidateJwtToken(users[0].Token); err != nil || !isValid {
			if !isValid {
				return errors.New("This account have been logged out already")
			} else {
				if !strings.Contains(err.Error(), "token is expired") {
					if err != nil {
						return err
					}
				}
			}
		}

		switch claim.UserType {
		case "STAFF":
			staffToken, staffExpiryTime, err := utils.GenerateJwtStaff(claim.UserId, claim.Username, claim.DutyLocationId)
			if err != nil {
				return err
			}

			updatedUser, err = UpdateUser(tx, claim.Username, claim.UserType, []string{"token_expiry_time", "token"}, utils.TimeInt64ToString(staffExpiryTime), staffToken)
			if err != nil {
				return err
			}

			return nil

		case "ADMIN":
			adminToken, adminExpiryTime, err := utils.GenerateJwtAdmin(claim.UserId, claim.Username, claim.Password)
			if err != nil {
				return err
			}

			updatedUser, err = UpdateUser(tx, claim.Username, claim.UserType, []string{"token_expiry_time", "token"}, utils.TimeInt64ToString(adminExpiryTime), adminToken)
			if err != nil {
				return err
			}
			return nil
		default:
			return errors.New("Unknown user type inside JWT token")
		}
		return nil
	})
	return updatedUser, err
}

func UsersToLoginResponse(user db_models.Users) (loginAuth.LoginResponseBody, error) {
	lastLogin, err := StringToResponseTime(user.LastLoginTime)
	if err != nil {
		return loginAuth.LoginResponseBody{}, err
	}
	expireTime, err := StringToResponseTime(user.TokenExpiryTime)
	if err != nil {
		return loginAuth.LoginResponseBody{}, err
	}
	return loginAuth.LoginResponseBody{ID: user.UserId, Username: user.Username, UserType: user.UserType, AuthToken: user.Token, LoginTime: lastLogin, TokenExpiryTime: expireTime, DutyLocationId: user.DutyLocationId}, nil
}

func UsersToUsersLogs(users []db_models.Users) ([]db_models.UsersLogs, error) {
	var usersLogsList []db_models.UsersLogs

	bJson, err := json.Marshal(users)
	if err != nil {
		return usersLogsList, err
	}
	err = json.Unmarshal(bJson, &usersLogsList)
	if err != nil {
		return usersLogsList, err
	}

	// for _, user := range users {
	// 	usersLogs := db_models.UsersLogs{Id: 0, UserId: user.UserId, Username: user.Username, Password: user.Password, UserType: user.UserType, Token: user.Token, DutyLocationId: user.DutyLocationId}
	// 	usersLogsList = append(usersLogsList, usersLogs)
	// }

	for _, usersLog := range usersLogsList {
		// usersLog.Id = 123
		if usersLog.TokenExpiryTime == "" {
			usersLog.TokenExpiryTime = utils.TimeInt64ToString(0)
		} else {
			time, err := StringToDatetime(usersLog.TokenExpiryTime)
			if err != nil {
				return usersLogsList, errors.New("Fail translate TokenExpiryTime")
			}
			usersLog.TokenExpiryTime = time
		}
		if usersLog.LastLoginTime == "" {
			usersLog.LastLoginTime = utils.TimeInt64ToString(0)
		} else {
			time, err := StringToDatetime(usersLog.LastLoginTime)
			if err != nil {
				return usersLogsList, errors.New("Fail translate LastLoginTime")
			}
			usersLog.LastLoginTime = time
		}
		if usersLog.LastLogoutTime == "" {
			usersLog.LastLogoutTime = utils.TimeInt64ToString(0)
		} else {
			time, err := StringToDatetime(usersLog.LastLogoutTime)
			if err != nil {
				return usersLogsList, errors.New("Fail translate LastLogoutTime")
			}
			usersLog.LastLogoutTime = time
		}
		if usersLog.CreateTime == "" {
			usersLog.CreateTime = utils.TimeInt64ToString(0)
		} else {
			time, err := StringToDatetime(usersLog.CreateTime)
			if err != nil {
				return usersLogsList, errors.New("Fail translate CreateTime")
			}
			usersLog.CreateTime = time
		}
		if usersLog.LastUpdateTime == "" {
			usersLog.LastUpdateTime = utils.TimeInt64ToString(0)
		} else {
			time, err := StringToDatetime(usersLog.LastUpdateTime)
			if err != nil {
				return usersLogsList, errors.New("Fail translate LastUpdateTime")
			}
			usersLog.LastUpdateTime = time
		}
	}

	return usersLogsList, nil
}
