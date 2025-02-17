package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"
	"fmt"
	"strconv"

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

func (c *BuildingController) GetBuilding(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var building = &[]models.BuildingModel{}

	isAuthenticated, err := c.buildingService.GetBuilding(ctx, building)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isAuthenticated {
		response.Message = config.GetMessageCode("INVALID_CREDENTIALS")
		ctx.JSON(401, response)
		return
	}

	response.Data = building
	ctx.JSON(200, response)
}

func (c *BuildingController) GetBuildingRoom(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	var buildingStruct = structs.BuildingID{
		ID: id,
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

func (c *BuildingController) CreateBuilding(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var building = &structs.NewBuilding{}

	if err := ctx.Bind(building); err != nil {
		// response.Message = config.GetMessageCode("INVALID_PARAMETER")
		// ctx.JSON(400, response)
		// return
		fmt.Println(err)
	}

	// if err := utils.Validate.Struct(building); err != nil {
	// 	response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
	// 	response.ValidateError = err.Error()
	// 	ctx.JSON(400, response)
	// 	return
	// }

	fmt.Println(building.Name)
	fmt.Println(building.Address)
	fmt.Println(building.Floors)
	fmt.Println(building.Managers)
	fmt.Println(building.Services)

	// if err := c.buildingService.CreateBuilding(ctx, building); err != nil {
	// 	response.Message = config.GetMessageCode("SYSTEM_ERROR")
	// 	ctx.JSON(500, response)
	// 	return
	// }
	ctx.JSON(400, response)
	// ctx.JSON(200, response)
}
