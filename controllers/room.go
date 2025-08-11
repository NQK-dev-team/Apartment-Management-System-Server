package controllers

import (
	"api/config"
	"api/services"
	"api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomService *services.RoomService
}

func NewRoomController() *RoomController {
	return &RoomController{
		roomService: services.NewRoomService(),
	}
}

func (c *RoomController) GetRoomList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	rooms := &[]structs.BuildingRoom{}

	if err := c.roomService.GetRoomList(ctx, rooms); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = rooms
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
