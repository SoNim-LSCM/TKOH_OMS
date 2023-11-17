package handlers

import (
	"encoding/json"

	// "github.com/SoNim-LSCM/TKOH_OMS/models"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models/systemStatus"
	"github.com/SoNim-LSCM/TKOH_OMS/models/wsTest"
	"github.com/SoNim-LSCM/TKOH_OMS/websocket"

	"github.com/gofiber/fiber/v2"
)

const AW2_RESPONSE string = `{
	"messageCode": "LOCATION_UPDATE",
	"userId": 1,
	"dutyLocationId" : 1,
	"dutyLocationName" : "5/F DSC"
}`

// @Summary		Test AW2 websocket response.
// @Description	Get the response of AW2 (Server notify the user which location selected).
// @Tags			Test
// @Param parameters body string false "AW2 response"
// @Produce		plain
// @Success		200	"OK"
// @Router			/testAW2 [post]
func HandleTestAW2(c *fiber.Ctx) error {

	// mqtt.PublishMqtt("direct/publish", []byte("packet scheduled message"))
	var response wsTest.ReportDutyLocationUpdateResponse
	if err := c.BodyParser(&response); errorHandler.CheckError(err, "Invalid/missing input: ") {
		err := json.Unmarshal([]byte(AW2_RESPONSE), &response)
		errorHandler.CheckError(err, "translate string to json in wsTest")
	}
	if err := websocket.SendBoardcastMessage(response); err != nil {
		return c.SendString(err.Error())
	}
	return c.SendString("OK")
}

const OW1_RESPONSE string = `{
    "messageCode": "ORDER_STATUS",
    "scheduleId": 1,
    "orderId": 1,
    "robotId":["AMR1"],
    "payloadId": "CART1",
    "orderStatus": "PROCESSING",
    "processingStatus": "ARRIVED_START_LOCATION"
}`

// @Summary		Test OW1 websocket response.
// @Description	Get the response of OW1 (Server notify any of created order status changed).
// @Tags			Test
// @Param parameters body string false "OW1 response"
// @Produce		plain
// @Success		200	"OK"
// @Router			/testOW1 [post]
func HandleTestOW1(c *fiber.Ctx) error {
	// get the processingStatus from the request body
	// processingStatus := c.Query("processingStatus")
	var response wsTest.ReportOrderStatusUpdateResponse
	if err := c.BodyParser(&response); errorHandler.CheckError(err, "Invalid/missing input: ") {
		err := json.Unmarshal([]byte(OW1_RESPONSE), &response)
		errorHandler.CheckError(err, "translate string to json in wsTest")
	}
	// err := json.Unmarshal([]byte(OW1_RESPONSE), &response)
	// fmt.Println(response)
	// response.ProcessingStatus = processingStatus
	// errorHandler.CheckError(err, "translate string to json in wsTest")
	if err := websocket.SendBoardcastMessage(response); err != nil {
		return c.SendString(err.Error())
	}
	return c.SendString("OK")
}

const MW1_RESPONSE string = `{
    "messageCode": "ROBOT_STATUS",
    "robotList": [
        {
            "robotId": "AMR1",
            "robotCoordatination": [255, 0],
            "robotPostion": [0.2, 0.0, 0.5 ],
            "robotOritenation": [0.01, 0.13, 0.0, 0.0],
            "robotState": "BUSY",
            "robotStatus": ["MOVE"],
            "batteryLevel": 89.5,
            "lastReportTime": "202310120800"
        }
    ]
}`

// @Summary		Test MW1 websocket response.
// @Description	Get the response of MW1 (Server report robot status and location (every 1s) ).
// @Tags			Test
// @Param parameters body string false "MW1 response"
// @Produce		plain
// @Success		200	"OK"
// @Router			/testHW1 [post]
func HandleTestMW1(c *fiber.Ctx) error {
	var response wsTest.ReportRobotStatusLocationResponse
	if err := c.BodyParser(&response); errorHandler.CheckError(err, "Invalid/missing input: ") {
		err := json.Unmarshal([]byte(MW1_RESPONSE), &response)
		errorHandler.CheckError(err, "translate string to json in wsTest")
	}
	if err := websocket.SendBoardcastMessage(response); err != nil {
		return c.SendString(err.Error())
	}
	return c.SendString("OK")
}

const SW1_RESPONSE string = `{
    "messageCode": "SYSTEM_STATUS",
    "systemState": "STOPPED",
    "systemStatus": ["LIFT_ALARM"]
}`

// @Summary		Test SW1 websocket response.
// @Description	Get the response of SW1 (Server report robot status and location (every 1s) ).
// @Tags			Test
// @Param parameters body string false "SW1 response"
// @Produce		plain
// @Success		200	"OK"
// @Router			/testSW1 [post]
func HandleTestSW1(c *fiber.Ctx) error {
	var response systemStatus.SystemStatusResponse
	if err := c.BodyParser(&response); errorHandler.CheckError(err, "Invalid/missing input: ") {
		err := json.Unmarshal([]byte(SW1_RESPONSE), &response)
		errorHandler.CheckError(err, "translate string to json in wsTest")
	}
	if err := websocket.SendBoardcastMessage(response); err != nil {
		return c.SendString(err.Error())
	}
	return c.SendString("OK")
}
