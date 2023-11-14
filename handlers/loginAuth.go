package handlers

import (
	"log"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"
	"github.com/SoNim-LSCM/TKOH_OMS/service"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"

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

	updatedUser, err := service.LoginStaff(request)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Login staff failed " + err.Error()))
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

	updatedUser, err := service.LoginAdmin(request)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Login staff failed " + err.Error()))
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
	_, err = service.Logout(claim, token)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Login failed: " + err.Error()))
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

	updatedUser, err := service.RenewToken(claim, token)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

	body, err := service.UsersToLoginResponse(updatedUser[0])
	if errorHandler.CheckError(err, "Translate from users to login response failed: ") {
		return c.Status(400).JSON(models.GetFailResponse("Translate from users to login response failed " + err.Error()))
	}

	response := loginAuth.LoginResponse{Header: header, Body: body}

	log.Printf("Renew token successful for user: %s\n", claim.Username)

	// return the API Response
	return c.Status(200).JSON(response)
}
