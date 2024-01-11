package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/constants/orderStatus"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
	ws_model "github.com/SoNim-LSCM/TKOH_OMS/models/websocket"
	"github.com/SoNim-LSCM/TKOH_OMS/service"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
	"github.com/SoNim-LSCM/TKOH_OMS/websocket"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get Delivery Order.
// @Description	Get the list of delivery order by order status which start/end at the staff's duty location.
// @Tags			Order Management
// @Accept			*/*
// @Param			orderStatus	query	[]string	false	"Order Status"
// @Param			scheduleId	query	int	false	"Schedule ID"
//
// @Produce		json
// @Success		200	{object} orderManagement.OrderListBody
// @Failure     400 {object} models.FailResponse
//
// @Router			/getDeliveryOrder [get]
// @Security Bearer
func HandleGetDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	// claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Get Delivery Order Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Get Delivery Order by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the orderStatus from the request body
	orders := []db_models.Orders{}
	statusString := c.Query("orderStatus")
	scheduleId := 0
	statusArray := strings.Split(statusString, ",")
	timeToIncludeCompletedOrders := time.Now().Add(time.Hour * -5).Format("2006-01-02 15:04:05")
	if statusString == "" || len(statusArray) == 0 {
		scheduleId, err = strconv.Atoi(c.Query("scheduleId"))
		if err != nil {
			return c.Status(400).JSON(models.GetFailResponse("Get Delivery Order Failed with Invalid/Missing Input", err.Error()))
		}
		log.Printf("Get Delivery Order with scheduleId: %s\n", fmt.Sprint(scheduleId))
		if claim.UserType == "ADMIN" {
			orders, err = service.FindOrders("schedule_id = ?", scheduleId)
		} else {
			orders, err = service.FindOrdersForFrontPage("schedule_id = ? AND (start_location_id = ? OR end_location_id = ?) AND NOT (order_status = 'COMPLETED' AND actual_arrival_time < ?)", claim.DutyLocationId, scheduleId, claim.DutyLocationId, claim.DutyLocationId, timeToIncludeCompletedOrders)
		}
	} else {
		log.Printf("Get Delivery Order with statusString: %s\n", statusString)
		if claim.UserType == "ADMIN" {
			orders, err = service.FindOrders("order_status IN ?", statusArray)
		} else {
			orders, err = service.FindOrdersForFrontPage("order_status IN ? AND (start_location_id = ? OR end_location_id = ?) AND NOT (order_status = 'COMPLETED' AND actual_arrival_time < ?)", claim.DutyLocationId, statusArray, claim.DutyLocationId, claim.DutyLocationId, timeToIncludeCompletedOrders)
		}
	}
	// err := service.FindRecords(&user, "orders", "order_status", statusArray, &orders)
	if errorHandler.CheckError(err, "Get Delivery Order Failed with Failed to Search Record") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to Search Record", err.Error()))
	}

	if len(orders) == 0 {
		log.Printf("Get Delivery Order Failed with Order not found\n")
		return c.Status(400).JSON(models.GetFailResponse("Order not found", ""))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	orderList, err := service.OrderListToOrderResponse(orders)
	if errorHandler.CheckError(err, "Get Delivery Order Failed with Data Transformation Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Data Transformation Failed", err.Error()))
	}

	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.GetDeliveryOrderResponse{Header: header, Body: body}
	// return the API Response
	log.Printf("Get Delivery Order Success for User: %s (%s)\n", claim.Username, claim.UserType)
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
	if errorHandler.CheckError(err, "Add Delivery Order Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Add Delivery Order by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the parameters from the request body
	var request dto.AddDeliveryOrderDTO
	var orderList orderManagement.OrderList
	// validate the request body
	if err := c.BodyParser(&request); errorHandler.CheckError(err, "Add Delivery Order Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}
	log.Printf("Get Delivery Order with Paramters: %s\n", request)

	orderList, err = service.AddOrders(append([]dto.AddDeliveryOrderDTO{}, request), claim.UserId, "ADHOC")
	if errorHandler.CheckError(err, "Add Delivery Order Failed with Add Orders Fail") {
		return c.Status(400).JSON(models.GetFailResponse("Add Orders Fail", err.Error()))
	}
	// err = service.InitOrderToRFMS(claim.UserId, orderList)
	// if errorHandler.CheckError(err, "Add Delivery Order Failed with Send Order to RFMS") {
	// 	return c.Status(400).JSON(models.GetFailResponse("Add Orders Fail", err.Error()))
	// }
	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.AddDeliveryOrderResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateOrderResponse(orderList))
	// return the API Response
	log.Printf("Add Delivery Order Success\n")
	return c.Status(200).JSON(response)
}

// @Summary		Get Routine Delivery Order.
// @Description	Get the list of routine delivery orders.
// @Tags			Order Management
// @Accept			*/*
//
// @Produce		json
// @Success		200	{object} orderManagement.OrderListBody
// @Failure     400 {object} models.FailResponse
//
// @Router			/getRoutineDeliveryOrder [get]
// @Security Bearer
func HandleGetRoutineDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	// claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Get Routine Delivery Order Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Get Routine Delivery Order by User: %s (%s)\n", claim.Username, claim.UserType)

	routines, err := service.FindRoutines("routine_id != ?", 0)
	// err := service.FindRecords(&user, "orders", "order_status", statusArray, &orders)
	if errorHandler.CheckError(err, "Get Routine Delivery Order Failed with Failed to Search Record") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to Search Record", err.Error()))
	}

	if len(routines) == 0 {
		log.Printf("Get Routine Delivery Order Failed with Routines not found")
		return c.Status(400).JSON(models.GetFailResponse("Routines not found", ""))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	routineResponse, err := service.RoutineListToRoutineResponse(routines)
	if errorHandler.CheckError(err, "Get Routine Delivery Order Failed with Data Transformation Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Data Transformation Failed", err.Error()))
	}

	body := orderManagement.RoutineOrderListBody{RoutineOrderList: routineResponse}
	response := orderManagement.GetRoutineDeliveryOrderResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateRoutineResponse(routineResponse))
	log.Printf("Get Routine Delivery Order Success for User: %s (%s)\n", claim.Username, claim.UserType)
	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Add Routine.
// @Description	Create Routine delivery order.
// @Tags			Order Management
// @Accept			json
//
// @Param parameters body dto.AddRoutineDTO true "Add Routine Parameters"
//
// @Produce		json
// @Success		200	{object} orderManagement.AddRoutineResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/addRoutineDeliveryOrder [post]
// @Security Bearer
func HandleAddRoutineDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Add Routine Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Get Delivery Order by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the parameters from the request body
	var request dto.AddRoutineDTO
	var routineOrderList orderManagement.RoutineOrderList
	// validate the request body
	if err := c.BodyParser(&request); errorHandler.CheckError(err, "Add Routine Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}
	log.Printf("Add Routine Order with Paramters: %s\n", request)

	routineOrderList, err = service.AddRoutines(request, claim.UserId)
	if errorHandler.CheckError(err, "Add Routine Failed with Add Routine Fail") {
		return c.Status(400).JSON(models.GetFailResponse("Add Routine Fail", err.Error()))
	}
	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.RoutineOrderListBody{RoutineOrderList: routineOrderList}
	response := orderManagement.AddRoutineResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateRoutineResponse(routineOrderList))
	log.Printf("Add Routine Success for User: %s (%s)\n", claim.Username, claim.UserType)
	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Update Routine Delivery Order.
// @Description	Update Routine Delivery Order .
// @Tags			Order Management
//
// @Accept			json
//
// @Param todo body dto.UpdateRoutineDeliveryOrderDTO true "Update Delivery Order Parameters"
//
// @Produce		json
// @Success		200	{object} orderManagement.UpdateRoutineDeliveryOrderResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/updateRoutineDeliveryOrder [post]
// @Security Bearer
func HandleUpdateRoutineDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Update Routine Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Update Routine by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the todo from the request body
	request := dto.UpdateRoutineDeliveryOrderDTO{}

	// validate the request body
	err = c.BodyParser(&request)
	if errorHandler.CheckError(err, "Update Routine Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}
	log.Printf("Update Routine with Paramters: %s\n", request)

	routineList, err := service.UpdateRoutineOrders(claim.UserId, request)
	if err != nil {
		log.Printf("Update Routine Failed with Update orders Failed")
		return c.Status(400).JSON(models.GetFailResponse("Update orders Failed", err.Error()))
	}
	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.RoutineOrderListBody{RoutineOrderList: routineList}
	response := orderManagement.UpdateRoutineDeliveryOrderResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateRoutineResponse(routineList))
	log.Printf("Update Routine Success: %s (%s)\n", claim.Username, claim.UserType)
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
	claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Trigger Order Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Trigger Order by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the orderIds and scheduleId from the request body
	orderIdsString := c.Query("orderIds")
	fmt.Printf("orderIds: %s\n", orderIdsString)
	scheduleIdString := c.Query("scheduleId")
	fmt.Printf("scheduleId: %s\n", scheduleIdString)

	log.Printf("Trigger Handling Order with Paramters: orderIds: %s, scheduleId: %s\n", orderIdsString, scheduleIdString)

	if (orderIdsString != "" && scheduleIdString != "") || (orderIdsString == "" && scheduleIdString == "") {
		log.Printf("Trigger Order Failed with Invalid/Missing Input")
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}

	var orderList orderManagement.OrderList

	if orderIdsString != "" {
		_, err := service.OrderIdsToIntArray(orderIdsString)
		if errorHandler.CheckError(err, "Trigger Order Failed with Fail formating orderIds") {
			return c.Status(400).JSON(models.GetFailResponse("Fail formating orderIds", err.Error()))
		}
		orderList, err = service.TriggerOrderOrderIds(orderIdsString)
		if errorHandler.CheckError(err, "Trigger Order Failed with Fail triggering order") {
			return c.Status(400).JSON(models.GetFailResponse("Fail triggering order", err.Error()))
		}
	} else {
		scheduleId, err := strconv.Atoi(scheduleIdString)
		if errorHandler.CheckError(err, "Trigger Order Failed with Fail formating scheduleId") {
			return c.Status(400).JSON(models.GetFailResponse("Fail formating scheduleId", err.Error()))
		}
		orderList, err = service.TriggerOrderScheduleId(scheduleId)
		if errorHandler.CheckError(err, "Trigger Order Failed with Fail triggering order") {
			return c.Status(400).JSON(models.GetFailResponse("Fail triggering order", err.Error()))
		}
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.TriggerHandlingOrderResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateOrderResponse(orderList))
	log.Printf("Trigger Order Success: %s (%s)\n", claim.Username, claim.UserType)
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
// @Param			processingStatus	query	string	false	"Processing Status"
//
// @Produce		json
// @Success		200	{object} orderManagement.UpdateDeliveryOrderResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/updateDeliveryOrder [post]
// @Security Bearer
func HandleUpdateDeliveryOrder(c *fiber.Ctx) error {
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Update Routine Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Update Routine by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the todo from the request body
	request := dto.UpdateDeliveryOrderDTO{}

	// validate the request body
	err = c.BodyParser(&request)
	if errorHandler.CheckError(err, "Update Routine Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}
	log.Printf("Update Delivery Order with Paramters: %s\n", request)

	orderList, err := service.UpdateOrders(claim.UserId, request)
	if errorHandler.CheckError(err, "Update Routine Failed with Update Orders Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Update Orders Failed", err.Error()))
	}

	// ----------------------------- overwriting fields -----------------------------------

	processingStatusString := c.Query("processingStatus")
	fmt.Printf("processingStatusString: %s\n", processingStatusString)
	if processingStatusString != "" {
		ordersss := []db_models.Orders{}
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			service.UpdateRecords(tx, &[]db_models.Orders{}, "orders", utils.CreateMap([]string{"processing_status"}, processingStatusString), "schedule_id = ? AND order_status = ?", request.ScheduleID, orderStatus.Created)
			return nil
		})
		if err != nil {
			return c.Status(400).JSON(models.GetFailResponse("overwrite processingStatus Failed (1)", err.Error()))
		}
		ordersss, err := service.FindOrders("schedule_id = ? AND order_status = ?", request.ScheduleID, orderStatus.Created)
		if err != nil {
			return c.Status(400).JSON(models.GetFailResponse("overwrite processingStatus Failed (2)", err.Error()))
		}
		orderList, err = service.OrderListToOrderResponse(ordersss)
		if err != nil {
			return c.Status(400).JSON(models.GetFailResponse("overwrite processingStatus Failed (3)", err.Error()))
		}
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.UpdateDeliveryOrderResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateOrderResponse(orderList))
	// return the API Response
	return c.Status(200).JSON(response)
}

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
	// verify token
	claim, _, err := utils.CtxToClaim(c)
	if errorHandler.CheckError(err, "Get Delivery Order Failed with Invalid Token") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Token", err.Error()))
	}
	log.Printf("Cancel Delivery Order by User: %s (%s)\n", claim.Username, claim.UserType)
	// get the todo from the request body
	request := new(dto.CancelDeliveryOrderDTO)

	// validate the request body
	err = c.BodyParser(&request)
	if errorHandler.CheckError(err, "Cancel Delivery Order Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid/Missing Input", err.Error()))
	}

	log.Printf("Cancel Order with Paramters: %s\n", request)

	orderList, err := service.CancelOrders(request.ScheduleID)
	if errorHandler.CheckError(err, "Cancel Delivery Order Failed with Cancel orders Failed") {
		return c.Status(400).JSON(models.GetFailResponse("Cancel orders Failed", err.Error()))
	}

	header := models.GetSuccessResponseHeader()
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.CancelDeliveryOrderResponse{Header: header, Body: body}
	websocket.SendBoardcastMessage(ws_model.GetUpdateOrderResponse(orderList))
	log.Printf("Cancel Delivery Order Success for User: %s (%s)\n", claim.Username, claim.UserType)
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
	log.Printf("Report Job Status\n")

	// get the todo from the request body
	request := dto.ReportJobStatusDTO{}

	// validate the request body
	err := c.BodyParser(&request)
	if errorHandler.CheckError(err, "Report Job Status Failed with Invalid/Missing Input") {
		return c.Status(400).JSON(fiber.Map{"Invalid/Missing Input": err.Error()})
	}
	log.Printf("Report Job Status with Paramters: %s\n", request)

	orderList, err := service.UpdateOrderFromRFMS(request)
	if errorHandler.CheckError(err, "Report Job Status Failed to Write Database") {
		return c.Status(400).JSON(fiber.Map{"Report Job Status Failed to Write Database": err.Error()})
	}

	response := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "Success"}
	// websocket.SendBoardcastMessage(ws_model.GetUpdateOrderResponse(orderList))
	websocket.SendBoardcastMessage(ws_model.GetUpdateOrderStatusResponse(orderList[0].OrderID, orderList[0].OrderStatus, request.PayloadID, orderList[0].ProcessingStatus, append([]string{}, request.RobotID), orderList[0].ScheduleID))
	log.Printf("Report Job Status Success\n")

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

	if errorHandler.CheckError(err, "Get Delivery Order Failed with Invalid Login Auth") {
		return c.Status(400).JSON(models.GetFailResponse("Invalid Login Auth", err.Error()))
	}

	log.Printf("Report System Status by User: %s\n", username)

	users, err := service.FindUsers("username", []string{username})
	// err := service.FindRecords(&user, "orders", "order_status", statusArray, &orders)
	if errorHandler.CheckError(err, "Get Delivery Order Failed with Failed to Search") {
		return c.Status(400).JSON(models.GetFailResponse("Failed to Search", err.Error()))
	}

	if len(users) == 0 {
		log.Printf("Get Delivery Order Failed with User Not Found\n")
		return c.Status(400).JSON(models.GetFailResponse("User Not Found", err.Error()))
	}
	// if errorHandler.CheckError(err, "Find user: "+username+" with type: RFMS in database") {
	// 	return c.Status(400).JSON(models.GetFailResponse(err.Error()))
	// }

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password))
	if errorHandler.CheckError(err, "Get Delivery Order Failed with Incorrect password") {
		return c.Status(400).JSON(models.GetFailResponse("Incorrect password", err.Error()))
	}

	response := orderManagement.ReportSystemStatusResponse{SystemState: "STOPPED", SystemStatus: []string{"LIFT_ALARM"}}
	log.Printf("Get Delivery Order Success with User: %s\n" + username)
	// return the API Response
	return c.Status(200).JSON(response)
}
