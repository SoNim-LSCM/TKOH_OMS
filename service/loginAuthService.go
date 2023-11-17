package service

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
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

		if request.DutyLocationId != users[0].DutyLocationId {
			return errors.New("Failed to search for staff")
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

		token, expiryTime, err := utils.GenerateJwtStaff(users[0].UserId, users[0].Username, users[0].DutyLocationId)
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
			log.Printf("Attempt to logout logged out account with username: %s user type: %s\n", claim.Username, claim.UserType)
			return errors.New("Account logged out already")
		} else if users[0].Token != token {
			log.Printf("Attempt to logout logged out account with incorrect token, username: %s user type: %s token: %s\n", claim.Username, claim.UserType, token)
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
			log.Printf("Attempt to logout logged out account with username: %s user type: %s\n", claim.Username, claim.UserType)
			return errors.New("Account logged out already")
		} else if users[0].Token != token {
			log.Printf("Attempt to logout logged out account with incorrect token, username: %s user type: %s token: %s\n", claim.Username, claim.UserType, token)
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
			log.Println("Unknown user type inside JWT token")
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

	for i, _ := range usersLogsList {
		// usersLogsList[i].Id = 123
		if usersLogsList[i].TokenExpiryTime == "" {
			usersLogsList[i].TokenExpiryTime = utils.TimeInt64ToString(0)
		}
		if usersLogsList[i].LastLoginTime == "" {
			usersLogsList[i].LastLoginTime = utils.TimeInt64ToString(0)
		}
		if usersLogsList[i].LastLogoutTime == "" {
			usersLogsList[i].LastLogoutTime = utils.TimeInt64ToString(0)
		}
		if usersLogsList[i].CreateTime == "" {
			usersLogsList[i].CreateTime = utils.TimeInt64ToString(0)
		}
		if usersLogsList[i].LastUpdateTime == "" {
			usersLogsList[i].LastUpdateTime = utils.TimeInt64ToString(0)
		}
	}

	return usersLogsList, nil
}
