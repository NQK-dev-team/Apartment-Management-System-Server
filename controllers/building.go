package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BuildingController struct {
	buildingService *services.BuildingService
	roomService     *services.RoomService
	contractService *services.ContractService
}

func NewBuildingController() *BuildingController {
	return &BuildingController{
		buildingService: services.NewBuildingService(),
		roomService:     services.NewRoomService(),
		contractService: services.NewContractService(),
	}
}

func (c *BuildingController) GetBuilding(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var building = &[]models.BuildingModel{}

	// getAll := ctx.Query("getAll") == "true"
	getAll := false

	isAuthenticated, err := c.buildingService.GetBuilding(ctx, building, getAll)
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
	response.Message = config.GetMessageCode("GET_SUCCESS")
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

	building.Name = strings.TrimSpace(building.Name)
	building.Address = strings.TrimSpace(building.Address)

	form, _ := ctx.MultipartForm()
	buildingImages := form.File["images[]"]
	building.Images = buildingImages

	for index, service := range building.Services {
		building.Services[index].Name = strings.TrimSpace(service.Name)
	}

	for index, room := range building.Rooms {
		building.Rooms[index].Description = strings.TrimSpace(room.Description)

		roomNoStr := strconv.Itoa(room.No)
		roomImages := form.File["roomImages["+roomNoStr+"]"]
		building.Rooms[index].Images = roomImages
	}

	if err := constants.Validate.Struct(building); err != nil {
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

	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *BuildingController) GetBuildingDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	building := &models.BuildingModel{}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	if permission := c.buildingService.CheckManagerPermission(ctx, id); !permission {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(403, response)
		return
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
	response.Message = config.GetMessageCode("GET_SUCCESS")
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

	response.Message = config.GetMessageCode("DELETE_SUCCESS")
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
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *BuildingController) GetBuildingRoom(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	rooms := &[]models.RoomModel{}

	isAuthenticated, err := c.buildingService.GetBuildingRoom(ctx, buildingID, rooms)
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

	response.Data = rooms
	response.Message = config.GetMessageCode("GET_SUCCESS")
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

	if permission := c.buildingService.CheckManagerPermission(ctx, building.ID); !permission {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(403, response)
		return
	}

	building.Name = strings.TrimSpace(building.Name)
	building.Address = strings.TrimSpace(building.Address)

	for index, service := range building.Services {
		building.Services[index].Name = strings.TrimSpace(service.Name)
	}

	for index, service := range building.NewServices {
		building.NewServices[index].Name = strings.TrimSpace(service.Name)
	}

	form, _ := ctx.MultipartForm()
	building.NewBuildingImages = form.File["newBuildingImages[]"]

	for index, room := range building.NewRooms {
		building.NewRooms[index].Description = strings.TrimSpace(room.Description)

		roomNoStr := strconv.Itoa(room.No)
		roomImages := form.File["newRoomImages["+roomNoStr+"]"]
		building.NewRooms[index].Images = roomImages
	}

	for index, room := range building.Rooms {
		building.Rooms[index].Description = strings.TrimSpace(room.Description)

		roomNoStr := strconv.Itoa(room.No)
		roomImages := form.File["newRoomImages["+roomNoStr+"]"]
		building.Rooms[index].NewImages = roomImages

		var oldRoomData models.RoomModel
		if err := c.roomService.GetRoomDetail(ctx, &oldRoomData, room.ID); err != nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.JSON(500, response)
			return
		}

		deletedImageCounter := 0

		for _, deletedImageID := range building.DeletedRoomImages {
			for _, oldRoomImage := range oldRoomData.Images {
				if deletedImageID == oldRoomImage.ID {
					deletedImageCounter++
				}
			}
		}

		building.Rooms[index].TotalImage = len(oldRoomData.Images) + len(roomImages) - deletedImageCounter
	}

	var oldBuildingData models.BuildingModel
	if err := c.buildingService.GetBuildingDetail(ctx, &oldBuildingData, building.ID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	building.TotalImage = len(oldBuildingData.Images) + len(building.NewBuildingImages) - len(building.DeletedBuildingImages)

	if err := constants.Validate.Struct(building); err != nil {
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

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *BuildingController) GetRoomDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	room := &structs.BuildingRoom{}

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		buildingID = 0
	}

	roomID, err := strconv.ParseInt(ctx.Param("roomID"), 10, 64)
	if err != nil {
		roomID = 0
	}

	if permission := c.buildingService.CheckManagerPermission(ctx, buildingID); !permission {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(403, response)
		return
	}

	buildingModel := &models.BuildingModel{}
	roomModel := &models.RoomModel{}
	contracts := &[]structs.Contract{}

	if err := c.buildingService.GetBuildingDetail(ctx, buildingModel, buildingID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if err := c.roomService.GetRoomDetail(ctx, roomModel, roomID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if err := c.contractService.GetContractByRoomIDAndBuildingID(ctx, contracts, roomID, buildingID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if roomModel.ID == 0 || buildingModel.ID == 0 || roomModel.BuildingID != buildingModel.ID || roomModel.ID != roomID || buildingModel.ID != buildingID {
		response.Message = config.GetMessageCode("DATA_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	room.ID = roomModel.ID
	room.CreatedAt = roomModel.CreatedAt
	room.CreatedBy = roomModel.CreatedBy
	room.UpdatedAt = roomModel.UpdatedAt
	room.UpdatedBy = roomModel.UpdatedBy
	room.DeletedAt = roomModel.DeletedAt
	room.DeletedBy = roomModel.DeletedBy
	room.No = roomModel.No
	room.Floor = roomModel.Floor
	room.Area = roomModel.Area
	room.Status = roomModel.Status
	room.Description = roomModel.Description
	room.BuildingID = roomModel.BuildingID
	room.BuildingName = buildingModel.Name
	room.Images = roomModel.Images
	room.BuildingAddress = buildingModel.Address
	room.Contracts = *contracts

	response.Data = room
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *BuildingController) GetRoomContract(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		buildingID = 0
	}

	roomID, err := strconv.ParseInt(ctx.Param("roomID"), 10, 64)
	if err != nil {
		roomID = 0
	}

	if permission := c.buildingService.CheckManagerPermission(ctx, buildingID); !permission {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(403, response)
		return
	}

	contracts := &[]structs.Contract{}

	if err := c.contractService.GetContractByRoomIDAndBuildingID(ctx, contracts, roomID, buildingID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = contracts
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}
