package controllers

import (
	"api/config"
	"api/services"
	"api/structs"
	"net/http"
	"strconv"

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

func (c *RoomController) GetRoomByID(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	roomID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		roomID = 0
	}

	room := &structs.BuildingRoom{}
	if err := c.roomService.GetRoomDetail2(ctx, room, roomID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if room.ID == 0 {
		response.Message = config.GetMessageCode("DATA_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response.Data = room
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *RoomController) GetRoomSupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	roomID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		roomID = 0
	}

	supportTickets := &[]structs.SupportTicket{}
	if err := c.roomService.GetRoomSupportTicket(ctx, supportTickets, roomID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = supportTickets
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
