package handlers

import (
	"fmt"
	"strconv"
	"strings"

	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
	"github.com/SoNim-LSCM/TKOH_OMS/service"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get Delivery Order.
// @Description	Get the list of delivery order by order status .
// @Tags			Order Management
// @Accept			*/*
// @Param			orderStatus	query	[]string	true	"Order Status"
//
// @Produce		json
// @Success		200	{object} orderManagement.OrderListBody
// @Failure     400 {object} models.FailResponse
//
// @Router			/getDeliveryOrder [get]
// @Security Bearer
func HandleGetDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	_, _, err := utils.CtxToClaim(c)
	// claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Invalid token: ") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token: " + err.Error()))
	}
	// get the orderStatus from the request body
	statusString := c.Query("orderStatus")
	statusArray := strings.Split(statusString, ",")
	if statusString == "" || len(statusArray) == 0 {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/missing input: " + err.Error()))
	}

	orders := []db_models.Orders{}
	// err := service.FindRecords(&user, "orders", "order_status", statusArray, &orders)
	if errorHandler.CheckError(service.FindRecords(&orders, "orders", "order_status IN ?", statusArray), "Failed to search") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to search: " + err.Error()))
	}

	if len(orders) == 0 {
		return c.Status(400).JSON(models.GetFailResponse("Order not found"))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	orderList, err := service.OrderListToOrderResponse(orders)
	if errorHandler.CheckError(err, "Translate from db_models.Orders to orderManagement.OrderList failed") {
		return c.Status(400).JSON(models.GetFailResponse("Translate from db_models.Orders to orderManagement.OrderList failed"))
	}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.GetDeliveryOrderResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Add Delivery Order.
// @Description	Create adhoc delivery order.
// @Tags			Order Management
// @Accept			json
//
// @Param parameters body dto.AddDeliveryOrderDTO true "Add Delivery Order Parameters"
//
// @Produce		json
// @Success		200	{object} orderManagement.AddDeliveryOrderResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/addDeliveryOrder [post]
// @Security Bearer
func HandleAddDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Invalid token: ") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token: " + err.Error()))
	}
	// get the parameters from the request body
	var request dto.AddDeliveryOrderDTO
	var orderList orderManagement.OrderList
	// validate the request body
	if err := c.BodyParser(&request); errorHandler.CheckError(err, "Invalid/missing input: ") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/missing input: " + err.Error()))
	}

	orderList, err = service.AddOrders(request, claim.UserId)
	if errorHandler.CheckError(err, "Add Orders Fail: ") {
		return c.Status(400).JSON(models.GetFailResponse("Add Orders Fail: " + err.Error()))
	}
	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.AddDeliveryOrderResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Trigger Delivery Order.
// @Description	Notify system users are ready to handle the current order.
// @Tags			Order Management
// @Accept			*/*
//
// @Param			orderIds	query		int	false	"Order IDs"
// @Param	scheduleId	query		int	false	"Schedule IDs"
//
// @Produce		json
// @Success		200	{object} orderManagement.TriggerHandlingOrderResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/triggerHandlingOrder [get]
// @Security Bearer
func HandleTriggerHandlingOrder(c *fiber.Ctx) error {
	// verify token
	_, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Invalid token: ") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token: " + err.Error()))
	}
	// get the orderIds and scheduleId from the request body
	orderIdsString := c.Query("orderIds")
	fmt.Printf("orderIds: %s\n", orderIdsString)
	scheduleIdString := c.Query("scheduleId")
	fmt.Printf("scheduleId: %s\n", scheduleIdString)

	if (orderIdsString != "" && scheduleIdString != "") || (orderIdsString == "" && scheduleIdString == "") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/missing input: "))
	}

	var orderList orderManagement.OrderList

	if orderIdsString != "" {
		orderIds, err := service.OrderIdsToIntArray(orderIdsString)
		if errorHandler.CheckError(err, "Failed when formating orderIds") {
			return c.Status(400).JSON(models.GetFailResponse("Failed when formating orderIds: " + err.Error()))
		}
		orderList, err = service.TriggerOrderOrderIds(orderIds)
		if errorHandler.CheckError(err, "Failed when triggering order") {
			return c.Status(400).JSON(models.GetFailResponse("Failed when triggering order: " + err.Error()))
		}
	} else {
		scheduleId, err := strconv.Atoi(scheduleIdString)
		if errorHandler.CheckError(err, "Failed when formating scheduleId") {
			return c.Status(400).JSON(models.GetFailResponse("Failed when formating scheduleId: " + err.Error()))
		}
		orderList, err = service.TriggerOrderScheduleId(scheduleId)
		if errorHandler.CheckError(err, "Failed when triggering order") {
			return c.Status(400).JSON(models.GetFailResponse("Failed when triggering order: " + err.Error()))
		}
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.TriggerHandlingOrderResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Update Delivery Order.
// @Description	Update Non Started Delivery Order .
// @Tags			Order Management
//
// @Accept			json
//
// @Param todo body dto.UpdateDeliveryOrderDTO true "Update Delivery Order Parameters"
//
// @Produce		json
// @Success		200	{object} orderManagement.UpdateDeliveryOrderResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/updateDeliveryOrder [post]
// @Security Bearer
func HandleUpdateDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	_, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Invalid token: ") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid token: " + err.Error()))
	}
	// get the todo from the request body
	request := dto.UpdateDeliveryOrderDTO{}

	// validate the request body
	err = c.BodyParser(&request)
	if errorHandler.CheckError(err, "Insufficient input paramters") {
		return c.Status(400).JSON(models.GetFailResponse("Insufficient input paramters: " + err.Error()))
	}

	orderList, err := service.UpdateOrders(request)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Update orders failed" + err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.UpdateDeliveryOrderResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}

const CancelDeliveryOrder string = `[
	{
		"scheduleId": 1,
		"orderId": 1,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "CANCELED",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectedStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectedDeliveryTime": "202310121000"
	}
]`

// @Summary		Cancel Delivery Order.
// @Description	Update Non Started Delivery Order .
// @Tags			Order Management
//
// @Accept			json
//
// @Param todo body dto.CancelDeliveryOrderDTO true "Cancel Delivery Parameters"
//
// @Produce		json
// @Success		200	{object} orderManagement.CancelDeliveryOrderResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/cancelDeliveryOrder [post]
// @Security Bearer
func HandleCancelDeliveryOrder(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(dto.CancelDeliveryOrderDTO)

	// validate the request body
	err := c.BodyParser(request)
	if errorHandler.CheckError(err, "Insufficient input paramters") {
		return c.Status(400).JSON(models.GetFailResponse("Insufficient input paramters" + err.Error()))
	}

	orderList, err := service.CancelOrders(request.ScheduleID)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Cancel orders failed" + err.Error()))
	}

	header := models.GetSuccessResponseHeader()
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.CancelDeliveryOrderResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Report Job Status.
// @Description	Receive the delivery job updated status.
// @Tags			Order Management
//
// @Accept			json
//
// @Param todo body dto.ReportJobStatusDTO true "Return Job Status Parameters"
//
// @Produce		json
// @Success		200	{object} models.ResponseHeader
// @Failure     400 {object} models.FailResponse
//
// @Router			/reportJobStatus [post]
func HandleReportJobStatus(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(dto.ReportJobStatusDTO)

	// validate the request body
	err := c.BodyParser(request)
	if errorHandler.CheckError(err, "Insufficient input parameters") {
		return c.Status(400).JSON(fiber.Map{"Insufficient input parameters": err.Error()})
	}

	response := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "Success"}

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Report System Status.
// @Description	Get current system status.
// @Tags			Order Management
//
// @Accept		*/*
//
// @Produce		json
// @Success		200	{object} orderManagement.ReportSystemStatusResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/reportSystemStatus [post]
func HandleReportSystemStatus(c *fiber.Ctx) error {

	username, password, err := utils.CtxToAuth(c)

	if errorHandler.CheckError(err, "Translate Ctx to Basic Auth String") {
		return c.Status(400).JSON(models.GetFailResponse(err.Error()))
	}

	user := db_models.Users{}
	err = service.FindRecords(&user, "users", "username", []string{username})

	if errorHandler.CheckError(err, "Find user: "+username+" with type: RFMS in database") {
		return c.Status(400).JSON(models.GetFailResponse(err.Error()))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errorHandler.CheckError(err, "Incorrect password") {
		return c.Status(400).JSON(models.GetFailResponse("Incorrect password: " + err.Error()))
	}

	response := orderManagement.ReportSystemStatusResponse{SystemState: "STOPPED", SystemStatus: []string{"LIFT_ALARM"}}

	// return the API Response
	return c.Status(200).JSON(response)
}
