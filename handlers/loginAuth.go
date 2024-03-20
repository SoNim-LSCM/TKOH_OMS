package handlers

import (
	"log"

	errorHandler "tkoh_oms/errors"
	"tkoh_oms/models"
	dto "tkoh_oms/models/DTO"
	"tkoh_oms/models/loginAuth"
	"tkoh_oms/service"
	"tkoh_oms/utils"

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
	err := c.BodyParser(&request)
	if errorHandler.CheckError(err, "Login Staff Insufficient input paramters") {
		return c.Status(400).JSON(models.GetFailResponse("Insufficient input paramters", err.Error()))
	}

	log.Printf("Login Staff with Username: %s, Duty Location ID: %d\n", request.Username, request.DutyLocationId)

	updatedUser, err := service.LoginStaff(request)
	if errorHandler.CheckError(err, "Login Staff Failed with Login Staff Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Login Staff Failed", err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

	body, err := service.UsersToLoginResponse(updatedUser[0])
	if errorHandler.CheckError(err, "Login Staff Failed with Data Transformation Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Data Transformation Failed", err.Error()))
	}

	response := loginAuth.LoginResponse{Header: header, Body: body}

	log.Printf("Login Staff Success for User: %s\n", request.Username)

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
	err := c.BodyParser(&request)
	if errorHandler.CheckError(err, "Login Admin Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}

	log.Printf("Login Admin Failed with Username: %s\n", request.Username)

	updatedUser, err := service.LoginAdmin(request)
	if errorHandler.CheckError(err, "Login Admin Failed with Login Admin Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Login Admin Failed", err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

	body, err := service.UsersToLoginResponse(updatedUser[0])
	if errorHandler.CheckError(err, "Login Admin Failed with Data Transformation Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Login Admin Failed with Data Transformation Failed", err.Error()))
	}

	response := loginAuth.LoginResponse{Header: header, Body: body}

	log.Printf("Login Successful for Admin User: %s\n", request.Username)

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
	claim, token, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Logout Failed with Invalid token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token", err.Error()))
	}
	log.Printf("Logout with User: %s\n", claim.Username)
	_, err = service.Logout(claim, token)
	if errorHandler.CheckError(err, "Logout Failed with Login Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Login Failed", err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	response := loginAuth.LogoutResponse{Header: header}

	log.Printf("Logout Success for User: %s\n", claim.Username)

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
	if errorHandler.CheckError(err, "Renew Token with Invalid token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token", err.Error()))
	}
	log.Printf("Renew Token by User: %s (%s)\n", claim.Username, claim.UserType)
	updatedUser, err := service.RenewToken(claim, token)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}

	body, err := service.UsersToLoginResponse(updatedUser[0])
	if errorHandler.CheckError(err, "Renew Token with Data Transformation Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Renew Token with Data Transformation Failed", err.Error()))
	}

	response := loginAuth.LoginResponse{Header: header, Body: body}

	log.Printf("Renew Token Success for User: %s (%s)\n", claim.Username, body.UserType)

	// return the API Response
	return c.Status(200).JSON(response)
}
