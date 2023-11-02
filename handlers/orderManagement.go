package handlers

import (
	"encoding/json"
	"fmt"
	"tkoh_oms/errors"
	"tkoh_oms/models"
	"tkoh_oms/models/orderManagement"

	"github.com/gofiber/fiber/v2"
)

const DeliveryOrderRequest string = `[
	{
		"scheduleId": 1,
		"orderId": 1,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "PROCESSING",
		"startTime": "202310120800",
		"endTime": "202310121000",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000",
		"processingStatus": "QUEUE_AT_START_BAY"
	},
	{
		"scheduleId": 1,
		"orderId": 2,
		"orderType": "DELIVERY_ONLY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "PROCESSING",
		"startTime": "202310120800",
		"endTime": "202310121000",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 3,
		"endLocationName": "LG/F Dirty Zone",
		"expectingDeliveryTime": "202310121000",
		"processingStatus": "QUEUE_AT_START_BAY"
	},
	{
		"scheduleId": 1,
		"orderId": 3,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "PROCESSING",
		"startTime": "202310120800",
		"endTime": "202310121000",
		"startLocationId": 2,
		"startLocationName": "LG/F Clean Zone",
		"expectingStartTime": "202310120830",
		"endLocationId": 1,
		"endLocationName": "5/F DSC",
		"expectingDeliveryTime": "202310121000",
		"processingStatus": "QUEUE_AT_START_BAY"
	},
	{
		"scheduleId": 1,
		"orderId": 4,
		"orderType": "DELIVERY_ONLY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "PROCESSING",
		"startTime": "202310120800",
		"endTime": "202310121000",
		"startLocationId": 3,
		"startLocationName": "LG/F Dirty Zone",
		"expectingStartTime": "202310120830",
		"endLocationId": 1,
		"endLocationName": "5/F DSC",
		"expectingDeliveryTime": "202310121000",
		"processingStatus": "QUEUE_AT_START_BAY"
	}
]`

// @Summary		Get Delivery Order.
// @Description	Get the list of delivery order by order status .
// @Tags			Order Management
// @Accept			*/*
//
//	@Produce		json
//	@Success		200	{object} orderManagement.OrderListBody
//
// @Router			/getDeliveryOrder [get]
func HandleGetDeliveryOrder(c *fiber.Ctx) error {
	// get the todo from the request body
	status := c.Query("status")
	fmt.Println(status)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	var orderList orderManagement.OrderList
	err := json.Unmarshal([]byte(DeliveryOrderRequest), &orderList)
	errors.CheckError(err, "translate string to json in orderManagement")
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.GetDeliveryOrderResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}

const AddDeliveryOrderRequest string = `[
	{
		"scheduleId": 1,
		"orderId": 1,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "ADHOC",
		"orderCreatedBy": 1,
		"orderStatus": "CREATED",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000"
	},
	{
		"scheduleId": 1,
		"orderId": 2,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "ADHOC",
		"orderCreatedBy": 1,
		"orderStatus": "CREATED",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000"
	}
]`

type AddDeliveryOrderDTO struct {
	OrderType             string `json:"orderType"`
	NumberOfAmrRequire    int    `json:"numberOfAmrRequire"`
	StartLocationID       int    `json:"startLocationId"`
	StartLocationName     string `json:"startLocationName"`
	ExpectingStartTime    string `json:"expectingStartTime"`
	EndLocationID         int    `json:"endLocationId"`
	EndLocationName       string `json:"endLocationName"`
	ExpectingDeliveryTime string `json:"expectingDeliveryTime"`
}

// @Summary		Add Delivery Order.
// @Description	Create adhoc delivery order.
// @Tags			Order Management
// @Accept			json
//
// @Param todo body AddDeliveryOrderDTO true "Add Delivery Order Parameters"
//
//	@Produce		json
//	@Success		200	{object} orderManagement.AddDeliveryOrderResponse
//
// @Router			/addDeliveryOrder [post]
func HandleAddDeliveryOrder(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(AddDeliveryOrderDTO)

	// validate the request body
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Bad Input"))
	}

	var orderList orderManagement.OrderList
	err = json.Unmarshal([]byte(AddDeliveryOrderRequest), &orderList)
	errors.CheckError(err, "translate string to json in orderManagement")
	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.AddDeliveryOrderResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}

const TriggerHandlingOrder string = `[
	{
		"scheduleId": 1,
		"orderId": 1,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "PROCESSING",
		"startTime": "202310120800",
		"endTime": "202310121000",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000",
		"processingStatus": "QUEUE_AT_START_BAY"
	}
]`

// @Summary		Trigger Delivery Order.
// @Description	Notify system users are ready to handle the current order.
// @Tags			Order Management
// @Accept			*/*
//
//	@Param			orderIds	query		int	false	"Order IDs"
//	@Param	scheduleId	query		int	false	"Schedule IDs"
//
//	@Produce		json
//	@Success		200	{object} orderManagement.TriggerHandlingOrderResponse
//
// @Router			/triggerHandlingOrder [get]
func HandleTriggerHandlingOrder(c *fiber.Ctx) error {
	// get the todo from the request body
	orderIds := c.Query("orderIds")
	fmt.Println(orderIds)
	scheduleId := c.Query("scheduleId")
	fmt.Println(scheduleId)

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	var orderList orderManagement.OrderList
	err := json.Unmarshal([]byte(TriggerHandlingOrder), &orderList)
	errors.CheckError(err, "translate string to json in orderManagement")
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.TriggerHandlingOrderResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}

const UpdateDeliveryOrder string = `[
	{
		"scheduleId": 1,
		"orderId": 1,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "SCHEDULED",
		"orderCreatedBy": 1,
		"orderStatus": "CREATED",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000"
	},
	{
		"scheduleId": 1,
		"orderId": 12,
		"orderType": "PICK_AND_DELIVERY",
		"orderCreatedType": "ADHOC",
		"orderCreatedBy": 1,
		"orderStatus": "CREATED",
		"startLocationId": 1,
		"startLocationName": "5/F DSC",
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000"
	}
]`

type UpdateDeliveryOrderDTO struct {
	ScheduleID            int    `json:"scheduleId"`
	NumberOfAmrRequire    int    `json:"numberOfAmrRequire"`
	StartLocationID       int    `json:"startLocationId"`
	StartLocationName     string `json:"startLocationName"`
	ExpectingStartTime    string `json:"expectingStartTime"`
	EndLocationID         int    `json:"endLocationId"`
	EndLocationName       string `json:"endLocationName"`
	ExpectingDeliveryTime string `json:"expectingDeliveryTime"`
}

// @Summary		Update Delivery Order.
// @Description	Update Non Started Delivery Order .
// @Tags			Order Management
//
//	@Accept			json
//
// @Param todo body UpdateDeliveryOrderDTO true "Update Delivery Order Parameters"
//
//	@Produce		json
//	@Success		200	{object} orderManagement.UpdateDeliveryOrderResponse
//
// @Router			/updateDeliveryOrder [get]
func HandleUpdateDeliveryOrder(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(UpdateDeliveryOrderDTO)

	// validate the request body
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	var orderList orderManagement.OrderList
	err = json.Unmarshal([]byte(UpdateDeliveryOrder), &orderList)
	errors.CheckError(err, "translate string to json in orderManagement")
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
		"expectingStartTime": "202310120830",
		"endLocationId": 2,
		"endLocationName": "LG/F Clean Zone",
		"expectingDeliveryTime": "202310121000"
	}
]`

type CancelDeliveryOrderDTO struct {
	ScheduleID int `json:"scheduleId"`
}

// @Summary		Cancel Delivery Order.
// @Description	Update Non Started Delivery Order .
// @Tags			Order Management
//
//	@Accept			json
//
// @Param todo body CancelDeliveryOrderDTO true "Cancel Delivery Parameters"
//
//	@Produce		json
//	@Success		200	{object} orderManagement.CancelDeliveryOrderResponse
//
// @Router			/cancelDeliveryOrder [post]
func HandleCancelDeliveryOrder(c *fiber.Ctx) error {
	// get the todo from the request body
	request := new(CancelDeliveryOrderDTO)

	// validate the request body
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	var orderList orderManagement.OrderList
	err = json.Unmarshal([]byte(CancelDeliveryOrder), &orderList)
	errors.CheckError(err, "translate string to json in orderManagement")
	body := orderManagement.OrderListBody{OrderList: orderList}
	response := orderManagement.CancelDeliveryOrderResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}
