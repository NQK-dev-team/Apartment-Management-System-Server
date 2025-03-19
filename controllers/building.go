package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"
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

func (c *BuildingController) CreateBuilding(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var building = &structs.NewBuilding{}

	if err := ctx.ShouldBind(building); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	form, _ := ctx.MultipartForm()
	buildingImages := form.File["images[]"]
	building.Images = buildingImages

	for index, room := range building.Rooms {
		roomNoStr := strconv.Itoa(room.No)
		roomImages := form.File["roomImages["+roomNoStr+"]"]
		building.Rooms[index].Images = roomImages
	}

	if err := utils.Validate.Struct(building); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.CreateBuilding(ctx, building); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *BuildingController) GetBuildingDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	building := &models.BuildingModel{}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	if err := c.buildingService.GetBuildingDetail(ctx, building, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if building.ID == 0 {
		response.Message = config.GetMessageCode("DATA_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	response.Data = building
	ctx.JSON(200, response)
}

func (c *BuildingController) DeleteBuilding(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	if err := c.buildingService.DeleteBuilding(ctx, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *BuildingController) GetBuildingSchedule(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	schedule := &[]models.ManagerScheduleModel{}

	isAuthenticated, err := c.buildingService.GetBuildingSchedule(ctx, buildingID, schedule)
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

	response.Data = schedule
	ctx.JSON(200, response)
}

func (c *BuildingController) UpdateBuilding(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var building = &structs.EditBuilding{}

	if err := ctx.ShouldBind(building); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if building.ID == 0 {
		buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

		if err != nil {
			response.Message = config.GetMessageCode("INVALID_PARAMETER")
			ctx.JSON(400, response)
			return
		}

		building.ID = buildingID
	}

	form, _ := ctx.MultipartForm()
	building.NewBuildingImages = form.File["newBuildingImages[]"]

	for index, room := range building.NewRooms {
		roomNoStr := strconv.Itoa(room.No)
		roomImages := form.File["newRoomImages["+roomNoStr+"]"]
		building.NewRooms[index].Images = roomImages
	}

	for index, room := range building.Rooms {
		roomNoStr := strconv.Itoa(room.No)
		roomImages := form.File["newRoomImages["+roomNoStr+"]"]
		building.Rooms[index].Images = roomImages
	}

	if err := utils.Validate.Struct(building); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.UpdateBuilding(ctx, building); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(500, response)
	// ctx.JSON(200, response)
}
