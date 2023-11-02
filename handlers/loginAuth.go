package handlers

import (
	"fmt"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	"github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginStaffDTO struct {
	Username       string `json:"username" `
	DutyLocationId int    `json:"dutyLocationId"`
}

//	@Summary		Login to OMS.
//	@Description	Login to OMS.
//	@Tags			Login Auth
//	@Accept			json
//
// @Param todo body LoginStaffDTO true "Login Parameters"
//
//		@Produce		json
//		@Success		200	{object} loginAuth.LoginResponse
//	 @Failure      	400 {object} models.FailResponse
//		@Router			/loginStaff [post]
func HandleLoginStaff(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(LoginStaffDTO)

	// validate the request body
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}

	token, expiryTime, err := utils.GenerateJWTStaff(request.Username, request.DutyLocationId)
	if err != nil {
		errors.CheckError(err, "login")
	}

	fmt.Println(token)
	fmt.Println(expiryTime)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := loginAuth.LoginResponseBody{ID: 1, Username: "user1", UserType: "STAFF", AuthToken: token, LoginDateTime: "20231010083000", TokenExpiryDateTime: "20231010180000", DutyLocationId: 1, DutyLocationName: "5/F DSC"}
	response := loginAuth.LoginResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}

type LoginAdminDTO struct {
	Username string `json:"username" `
	Password string `json:"password"`
}

//	@Summary		Login to OMS.
//	@Description	Login to OMS.
//	@Tags			Login Auth
//	@Accept			json
//
// @Param todo body LoginAdminDTO true "Login Parameters"
//
//		@Produce		json
//		@Success		200	{object} loginAuth.LoginResponse
//	 @Failure      	400 {object} models.FailResponse
//		@Router			/loginAdmin [post]
func HandleLoginAdmin(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(LoginAdminDTO)

	// validate the request body
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}
	fmt.Println(request.Username)
	userInfo, err := database.FindUser(request.Username, request.Password, "ADMIN")
	if err != nil || userInfo == nil {
		return c.Status(400).JSON(fiber.Map{"failed to search for": err.Error()})
	}

	token, expiryTime, err := utils.GenerateJWTAdmin(request.Username, request.Password)
	errors.CheckError(err, "generate token")

	database.UpdateUser(userInfo[0], token, expiryTime)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := loginAuth.LoginResponseBody{ID: 1, Username: "user1", UserType: "ADMIN", AuthToken: token, LoginDateTime: "20231010083000", TokenExpiryDateTime: "20231010180000", DutyLocationId: 1, DutyLocationName: "5/F DSC"}
	response := loginAuth.LoginResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Logout from OMS.
// @Description	Logout from OMS.
// @Tags			Login Auth
// @Accept			json
// @Produce		json
// @Success		200	{object} loginAuth.LogoutResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/logout [get]
// @Security Bearer
func HandleLogout(c *fiber.Ctx) error {

	user, err := utils.ParseToken(c)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse(err.Error()))
	}
	fmt.Println(user.UserType)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	response := loginAuth.LogoutResponse{Header: header}

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Websocket Connection.
// @Description	Using Valid Token to renew token before expired
// @Tags			Login Auth
// @Accept			*/*
//
//		@Produce		json
//		@Success		200	{object} loginAuth.LogoutResponse
//	 @Failure      	400 {object} models.FailResponse
//
// @Router			/renewToken [get]
func HandleRenewToken(c *fiber.Ctx) error {

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := loginAuth.RenewTokenResponseBody{ID: 1, Username: "user1", AuthToken: "GENERATED_TOKEN", LoginDateTime: "20231010083000", TokenExpiryDateTime: "20231010180000"}
	response := loginAuth.RenewTokenResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}
