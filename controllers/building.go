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

func (c *BuildingController) GetBuildingService(ctx *gin.Context) {
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

	var service = &[]models.BuildingServiceModel{}

	if err := c.buildingService.GetBuildingService(ctx, buildingStruct.ID, service); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = service
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

func (c *BuildingController) DeleteRooms(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	idList := &structs.IDList{}

	if err := ctx.ShouldBind(idList); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(idList); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.DeleteRooms(ctx, buildingID, idList.IDs); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *BuildingController) DeleteServices(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	// buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	// if err != nil {
	// 	response.Message = config.GetMessageCode("INVALID_PARAMETER")
	// 	ctx.JSON(400, response)
	// 	return
	// }

	idList := &structs.IDList{}

	if err := ctx.ShouldBind(idList); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(idList); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.DeleteServices(ctx, idList.IDs); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *BuildingController) AddService(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	newService := &structs.Service{}

	if err := ctx.ShouldBind(newService); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(newService); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.AddService(ctx, &models.BuildingServiceModel{
		BuildingID: buildingID,
		Name:       newService.Name,
		Price:      newService.Price,
	}); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *BuildingController) EditService(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	// buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	// if err != nil {
	// 	response.Message = config.GetMessageCode("INVALID_PARAMETER")
	// 	ctx.JSON(400, response)
	// 	return
	// }

	serviceID, err := strconv.ParseInt(ctx.Param("serviceID"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	service := &structs.Service{}

	if err := ctx.ShouldBind(service); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(service); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.EditService(ctx, &models.BuildingServiceModel{
		// BuildingID: buildingID,
		Name:  service.Name,
		Price: service.Price,
		DefaultModel: models.DefaultModel{
			ID: serviceID,
		},
	}); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *BuildingController) AddRoom(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	newRoom := &structs.NewRoom{}

	if err := ctx.ShouldBind(newRoom); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	form, _ := ctx.MultipartForm()
	newRoom.Images = form.File["images[]"]

	if err := utils.Validate.Struct(newRoom); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.buildingService.AddRoom(ctx, buildingID, newRoom); err != nil {
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

	if err := c.buildingService.GetBuildingSchedule(ctx, buildingID, schedule); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = schedule
	ctx.JSON(200, response)
}
