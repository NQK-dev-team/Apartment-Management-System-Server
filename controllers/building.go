package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"

	"github.com/gin-gonic/gin"
)

type BuildingController struct {
	buildingService *services.BuildingService
}

func NewBuildingController() *BuildingController {
	return &BuildingController{
		buildingService: services.NewBuildingService(),
	}
}

func (c *BuildingController) GetBuildings(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var buildings = &[]models.BuildingModel{}

	if err := c.buildingService.GetBuilding(ctx, buildings); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = buildings
	ctx.JSON(200, response)
}

func (c *BuildingController) GetBuildingRoom(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	var buildingStruct = structs.BuildingID{
		ID: ctx.Param("id"),
	}

	if err := utils.Validate.Struct(&buildingStruct); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	var room = &[]models.RoomModel{}

	if err := c.buildingService.GetBuildingRoom(ctx, buildingStruct.ID, room); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = room
	ctx.JSON(200, response)
}
