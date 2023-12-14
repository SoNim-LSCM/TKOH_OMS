package handlers

import (
	"encoding/json"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	"github.com/SoNim-LSCM/TKOH_OMS/models/mapHandling"
	"github.com/SoNim-LSCM/TKOH_OMS/service"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get Floor Plan.
// @Description	Get UI Floor Plan.
// @Tags			Map Handling
// @Accept			*/*
// @Produce		json
// @Success		200	{object} mapHandling.GetFloorPlanResponse
// @Failure     400 {object} models.FailResponse
//
// @Router			/getFloorPlan [get]
func HandleGetFloorPlan(c *fiber.Ctx) error {

	floorPlan, err := service.GetFloorPlan()
	if errorHandler.CheckError(err, "Get Floor Plan from Database") {
		return c.Status(400).JSON(models.GetFailResponse("Get Floor Plan from Database", err.Error()))
	}
	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	mapList, err := service.FloorPlanToMapList(floorPlan)
	if errorHandler.CheckError(err, "Translate string to json in mapHandling") {
		return c.Status(400).JSON(models.GetFailResponse("translate string to json in mapHandling", err.Error()))
	}
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
// @Failure     400 {object} models.FailResponse
//
// @Router			/getDutyRooms [get]
func HandleGetDutyRooms(c *fiber.Ctx) error {

	err := service.GetLocationFromRFMS()
	if errorHandler.CheckError(err, "Get Location from RFMS") {
		return c.Status(400).JSON(models.GetFailResponse("Get Location from RFMS", err.Error()))
	}

	var locationList mapHandling.LocationList
	mainInterface, err := service.FindAllDutyRooms()
	if errorHandler.CheckError(err, "Duty Rooms Not Found") {
		return c.Status(400).JSON(models.GetFailResponse("Duty Rooms Not Found", err.Error()))
	}
	jsonString, err := json.Marshal(mainInterface)
	if errorHandler.CheckError(err, "Translate struct to json string in mapHandling") {
		return c.Status(400).JSON(models.GetFailResponse("Translate struct to json string in mapHandling", err.Error()))
	}
	err = json.Unmarshal(jsonString, &locationList)
	if errorHandler.CheckError(err, "Translate json to struct in mapHandling") {
		return c.Status(400).JSON(models.GetFailResponse("Translate string to json in mapHandling", err.Error()))
	}

	header := models.ResponseHeader{ResponseCode: 200, ResponseMessage: "SUCCESS"}
	body := mapHandling.LocationListBody{LocationList: locationList}
	response := mapHandling.GetDutyRoomsResponse{Header: header, Body: body}
	// return the API Response
	return c.Status(200).JSON(response)
}
