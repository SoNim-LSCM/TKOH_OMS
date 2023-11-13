package handlers

import (
	"log"
	"strings"

	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"
	"github.com/SoNim-LSCM/TKOH_OMS/service"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Login to OMS.
// @Description	Login to OMS.
// @Tags			Login Auth
// @Accept			json
//
// @Param todo body dto.LoginStaffDTO true "Login Parameters"
//
// @Produce		json
// @Success		200	{object} loginAuth.LoginResponse
// @Failure      	400 {object} models.FailResponse
// @Router			/loginStaff [post]
func HandleLoginStaff(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(dto.LoginStaffDTO)

	// validate the request body
	err := c.BodyParser(request)
	if errorHandler.CheckError(err, "HandleLoginStaff Insufficient input paramters") {
		return c.Status(400).JSON(models.GetFailResponse("Insufficient input paramters: " + err.Error()))
	}

	log.Printf("HandleLoginStaff with username: %s dutyLocationId: %d\n", request.Username, request.DutyLocationId)

	user := db_models.Users{}
	if errorHandler.CheckError(service.FindRecords(&user, "users", &db_models.Users{Username: request.Username, UserType: "STAFF"}), "HandleLoginStaff failed to search") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to search for record"))
	}

	if user.UserId == 0 {
		return c.Status(400).JSON(models.GetFailResponse("Order not found"))
	}

	if request.DutyLocationId != user.DutyLocationId {
		return c.Status(400).JSON(models.GetFailResponse("Failed to search for staff"))
	}

	if isValid, err := utils.ValidateJwtToken(user.Token); err != nil || isValid {
		if isValid {
			return c.Status(400).JSON(models.GetFailResponse("this account is already logged in by someone"))
		} else {
			if !strings.Contains(err.Error(), "token is expired") {
				errorHandler.CheckError(err, "HandleLoginStaff failed to validate token")
				return c.Status(400).JSON(models.GetFailResponse("failed to validate token: " + err.Error()))
			}
		}
	}

	token, expiryTime, err := utils.GenerateJwtStaff(user.UserId, user.Username, user.DutyLocationId)
	if errorHandler.CheckError(err, "HandleLoginStaff generate token fail") {
		return c.Status(400).JSON(models.GetFailResponse("generate token fail: " + err.Error()))
	}

	// updatedUser := []db_models.Users{}
	// timeNow := utils.GetTimeNowString()
	// updateFields := []string{"token_expiry_datetime", "last_login_datetime", "token", "duty_location_id"}
	// updateMap := utils.CreateMap(updateFields, utils.TimeInt64ToString(expiryTime), timeNow, token, request.DutyLocationId)
	updatedUser, err := service.UpdateUser(user.Username, user.UserType, []string{"token_expiry_datetime", "last_login_datetime", "token", "duty_location_id"}, utils.TimeInt64ToString(expiryTime), utils.GetTimeNowString(), token, request.DutyLocationId)
	// err = service.UpdateRecords(&updatedUser, "users", updateMap, "username = ? AND user_type = ?", user.Username, user.UserType)

	if errorHandler.CheckError(err, "Update user fail: ") {
		return c.Status(400).JSON(models.GetFailResponse("Update user fail " + err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

	body, err := service.UsersToLoginResponse(updatedUser[0])
	if errorHandler.CheckError(err, "Translate from users to login response failed: ") {
		return c.Status(400).JSON(models.GetFailResponse("Translate from users to login response failed " + err.Error()))
	}

	response := loginAuth.LoginResponse{Header: header, Body: body}

	log.Printf("HandleLoginStaff login successful for staff user: %s\n", request.Username)

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Login to OMS.
// @Description	Login to OMS.
// @Tags		Login Auth
// @Accept		json
//
// @Param todo body dto.LoginAdminDTO true "Login Parameters"
//
// @Produce		json
// @Success		200	{object} loginAuth.LoginResponse
// @Failure      400 {object} models.FailResponse
// @Router		/loginAdmin [post]
func HandleLoginAdmin(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(dto.LoginAdminDTO)

	// validate the request body
	err := c.BodyParser(request)
	if errorHandler.CheckError(err, "HandleLoginAdmin Insufficient input paramters") {
		return c.Status(400).JSON(models.GetFailResponse("Insufficient input paramters: " + err.Error()))
	}

	log.Printf("HandleLoginAdmin with username: %s password: %s\n", request.Username, request.Password)

	user := db_models.Users{}
	if errorHandler.CheckError(service.FindRecords(&user, "users", &db_models.Users{Username: request.Username, UserType: "ADMIN"}), "HandleLoginAdmin failed to search") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to search for record"))
	}

	if isValid, err := utils.ValidateJwtToken(user.Token); err != nil || isValid {
		if isValid {
			return c.Status(400).JSON(models.GetFailResponse("this account is already logged in by someone"))
		} else {
			if !strings.Contains(err.Error(), "token is expired") {
				errorHandler.CheckError(err, "HandleLoginAdmin failed to validate token")
				return c.Status(400).JSON(models.GetFailResponse("failed to validate token: " + err.Error()))
			}
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if errorHandler.CheckError(err, "HandleLoginAdmin incorrect password") {
		return c.Status(400).JSON(models.GetFailResponse("incorrect password: " + err.Error()))
	}

	token, expiryTime, err := utils.GenerateJwtAdmin(user.UserId, user.Username, user.Password)
	if errorHandler.CheckError(err, "HandleLoginAdmin generate token fail") {
		return c.Status(400).JSON(models.GetFailResponse("generate token fail: " + err.Error()))
	}

	updatedUser, err := service.UpdateUser(user.Username, user.UserType, []string{"token_expiry_datetime", "last_login_datetime", "token"}, utils.TimeInt64ToString(expiryTime), utils.GetTimeNowString(), token)
	if errorHandler.CheckError(err, "Update user fail: ") {
		return c.Status(400).JSON(models.GetFailResponse("Update user fail " + err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

	body, err := service.UsersToLoginResponse(updatedUser[0])
	if errorHandler.CheckError(err, "Translate from users to login response failed: ") {
		return c.Status(400).JSON(models.GetFailResponse("Translate from users to login response failed " + err.Error()))
	}

	response := loginAuth.LoginResponse{Header: header, Body: body}

	log.Printf("Login successful for admin user: %s\n", request.Username)

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Logout from OMS.
// @Description	Logout from OMS.
// @Tags		Login Auth
// @Accept		json
// @Produce		json
// @Success		200	{object} loginAuth.LogoutResponse
// @Failure     400 {object} models.FailResponse
//
// @Router		/logout [get]
// @Security Bearer
func HandleLogout(c *fiber.Ctx) error {
	log.Printf("mysql query: HandleLogout start\n")
	claim, token, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Invalid token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token: " + err.Error()))
	}
	log.Printf("mysql query: HandleLogout: %s, %s\n", claim.Username, claim.UserType)
	user := db_models.Users{}
	if errorHandler.CheckError(service.FindRecords(&user, "users", &db_models.Users{Username: claim.Username, UserType: claim.UserType}), "Find user: ") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to search for record"))
	} else if user.Token == "" {
		log.Printf("Attempt to logout logged out account with username: %s user type: %s\n", claim.Username, claim.UserType)
		return c.Status(400).JSON(models.GetFailResponse("Account logged out already"))
	} else if user.Token != token {
		log.Printf("Attempt to logout logged out account with incorrect token, username: %s user type: %s token: %s\n", claim.Username, claim.UserType, token)
		return c.Status(400).JSON(models.GetFailResponse("Incorrect token"))
	}

	_, err = service.UpdateUser(user.Username, user.UserType, []string{"last_logout_datetime", "token"}, utils.GetTimeNowString(), "")
	if errorHandler.CheckError(err, "Fail update database") {
		return c.Status(400).JSON(models.GetFailResponse("Fail update database: " + err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	response := loginAuth.LogoutResponse{Header: header}

	log.Printf("Logout successful for user: %s\n", claim.Username)

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Renew JWT Token.
// @Description	Using Valid Token to renew token before expired
// @Tags		Login Auth
// @Accept		*/*
//
// @Produce		json
// @Success		200	{object} loginAuth.LoginResponse
// @Success		200	{object} loginAuth.LoginResponse
// @Failure     400 {object} models.FailResponse
//
// @Router		/renewToken [get]
// @Security Bearer
func HandleRenewToken(c *fiber.Ctx) error {

	claim, token, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Invalid token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token: " + err.Error()))
	}

	user := db_models.Users{}
	if errorHandler.CheckError(service.FindRecords(&user, "users", &db_models.Users{Username: claim.Username, UserType: claim.UserType}), "Find user: ") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to search for record"))
	} else if user.Token == "" {
		log.Printf("Attempt to logout logged out account with username: %s user type: %s\n", claim.Username, claim.UserType)
		return c.Status(400).JSON(models.GetFailResponse("Account logged out already"))
	} else if user.Token != token {
		log.Printf("Attempt to logout logged out account with incorrect token, username: %s user type: %s token: %s\n", claim.Username, claim.UserType, token)
		return c.Status(400).JSON(models.GetFailResponse("Incorrect token"))
	}

	if isValid, err := utils.ValidateJwtToken(user.Token); err != nil || !isValid {
		if !isValid {
			return c.Status(400).JSON(models.GetFailResponse("This account have been logged out already"))
		} else {
			if !strings.Contains(err.Error(), "token is expired") {
				errorHandler.CheckError(err, "Failed to validate token")
				return c.Status(400).JSON(models.GetFailResponse("Failed to validate token: " + err.Error()))
			}
		}
	}

	switch claim.UserType {
	case "STAFF":
		staffToken, staffExpiryTime, err := utils.GenerateJwtStaff(claim.UserId, claim.Username, claim.DutyLocationId)
		if errorHandler.CheckError(err, "Generate token fail") {
			return c.Status(400).JSON(models.GetFailResponse("Generate token fail: " + err.Error()))
		}

		updatedUser, err := service.UpdateUser(claim.Username, claim.UserType, []string{"token_expiry_datetime", "token"}, utils.TimeInt64ToString(staffExpiryTime), staffToken)
		if errorHandler.CheckError(err, "Update user fail: ") {
			return c.Status(400).JSON(models.GetFailResponse("Update user fail " + err.Error()))
		}

		header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

		body, err := service.UsersToLoginResponse(updatedUser[0])
		if errorHandler.CheckError(err, "Translate from users to login response failed: ") {
			return c.Status(400).JSON(models.GetFailResponse("Translate from users to login response failed " + err.Error()))
		}

		response := loginAuth.LoginResponse{Header: header, Body: body}

		log.Printf("Renew token successful for staff user: %s\n", claim.Username)

		// return the API Response
		return c.Status(200).JSON(response)

	case "ADMIN":
		adminToken, adminExpiryTime, err := utils.GenerateJwtAdmin(claim.UserId, claim.Username, claim.Password)
		if errorHandler.CheckError(err, "Generate token fail") {
			return c.Status(400).JSON(models.GetFailResponse("Generate token fail: " + err.Error()))
		}

		updatedUser, err := service.UpdateUser(claim.Username, claim.UserType, []string{"token_expiry_datetime", "token"}, utils.TimeInt64ToString(adminExpiryTime), adminToken)
		if errorHandler.CheckError(err, "Update user fail: ") {
			return c.Status(400).JSON(models.GetFailResponse("Update user fail " + err.Error()))
		}
		header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

		body, err := service.UsersToLoginResponse(updatedUser[0])
		if errorHandler.CheckError(err, "Translate from users to login response failed: ") {
			return c.Status(400).JSON(models.GetFailResponse("Translate from users to login response failed " + err.Error()))
		}

		response := loginAuth.LoginResponse{Header: header, Body: body}

		log.Printf("Renew token successful for admin user: %s\n", claim.Username)

		// return the API Response
		return c.Status(200).JSON(response)
	default:
		log.Println("Unknown user type inside JWT token")
		return c.Status(400).JSON(models.GetFailResponse("Unknown user type inside JWT token"))
	}
}
