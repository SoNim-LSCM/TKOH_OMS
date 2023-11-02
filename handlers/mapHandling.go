package handlers

import (
	"encoding/json"
	"tkoh_oms/database"
	"tkoh_oms/errors"
	"tkoh_oms/models"
	"tkoh_oms/models/mapHandling"

	"github.com/gofiber/fiber/v2"
)

const GET_FLOOR_PLAN_RESPONSE string = `[
	{
		"floorId": 1,
		"floorName": "5/F",
		"floorImage": "base64 image"
	},
	{
		"floorId": 2,
		"floorName": "LG/F",
		"floorImage": "base64 image"
	}
]
`

// @Summary		Get Floor Plan.
// @Description	Get UI Floor Plan.
// @Tags			Map Handling
// @Accept			*/*
// @Produce		json
// @Success		200	{object} mapHandling.GetFloorPlanResponse
// @Router			/getFloorPlan [get]
func HandleGetFloorPlan(c *fiber.Ctx) error {

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	var mapList mapHandling.MapList
	err := json.Unmarshal([]byte(GET_FLOOR_PLAN_RESPONSE), &mapList)
	errors.CheckError(err, "translate string to json in mapHandling")
	body := mapHandling.MapListBody{MapList: mapList}
	response := mapHandling.GetFloorPlanResponse{Header: header, Body: body}

	// return the API Response
	return c.Status(200).JSON(response)
}

// @Summary		Get Duty Rooms.
// @Description	Get the list of location.
// @Tags			Map Handling
// @Accept			*/*
// @Produce		json
// @Success		200	{object} mapHandling.GetDutyRoomsResponse
// @Router			/getDutyRooms [get]
func HandleGetDutyRooms(c *fiber.Ctx) error {

	var locationList mapHandling.LocationList
	mainInterface, err := database.FindAllDutyRooms()
	if err != nil {
		return c.Status(400).JSON(models.GetFailResponse("Duty Rooms Not Found"))
	}
	jsonString, err := json.Marshal(mainInterface)
	errors.CheckError(err, "translate struct to json string in mapHandling")
	err = json.Unmarshal(jsonString, &locationList)
	errors.CheckError(err, "translate string to json in mapHandling")

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := mapHandling.LocationListBody{LocationList: locationList}
	response := mapHandling.GetDutyRoomsResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}
