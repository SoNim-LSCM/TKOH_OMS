package handlers

import (
	"encoding/json"
	// "github.com/SoNim-LSCM/TKOH_OMS/models"

	"github.com/SoNim-LSCM/TKOH_OMS/errors"
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
// @Accept			*/*
// @Produce		plain
// @Success		200	"OK"
// @Router			/testAW2 [get]
func HandleTestAW2(c *fiber.Ctx) error {

	// mqtt.PublishMqtt("direct/publish", []byte("packet scheduled message"))
	var response wsTest.ReportDutyLocationUpdateResponse
	err := json.Unmarshal([]byte(AW2_RESPONSE), &response)
	errors.CheckError(err, "translate string to json in wsTest")
	websocket.SendMessage(response)
	return c.SendString("OK")
}

const OW1_RESPONSE string = `{
    "messageCode": "ORDER_ARRIVED",
    "scheduleId": 1,
    "orderId": 1,
    "robotId":["AMR1"],
    "orderStatus": "PROCESSING",
    "processingStatus": "ARRIVED_START_LOCATION"
}`

// @Summary		Test OW2 websocket response.
// @Description	Get the response of OW2 (Server notify any of created order status changed).
// @Tags			Test
// @Accept			*/*
// @Produce		plain
// @Success		200	"OK"
// @Router			/testOW1 [get]
func HandleTestOW1(c *fiber.Ctx) error {
	var response wsTest.ReportOrderStatusUpdateResponse
	err := json.Unmarshal([]byte(OW1_RESPONSE), &response)
	errors.CheckError(err, "translate string to json in wsTest")
	websocket.SendMessage(response)
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
// @Accept			*/*
// @Produce		plain
// @Success		200	"OK"
// @Router			/testHW1 [get]
func HandleTestMW1(c *fiber.Ctx) error {
	var response wsTest.ReportRobotStatusLocationResponse
	err := json.Unmarshal([]byte(MW1_RESPONSE), &response)
	errors.CheckError(err, "translate string to json in wsTest")
	websocket.SendMessage(response)
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
// @Accept			*/*
// @Produce		plain
// @Success		200	"OK"
// @Router			/testSW1 [get]
func HandleTestSW1(c *fiber.Ctx) error {
	var response systemStatus.SystemStatusResponse
	err := json.Unmarshal([]byte(SW1_RESPONSE), &response)
	errors.CheckError(err, "translate string to json in wsTest")
	websocket.SendMessage(response)
	return c.SendString("OK")
}
