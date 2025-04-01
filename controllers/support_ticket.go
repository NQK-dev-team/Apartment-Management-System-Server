package controllers

import (
	"api/config"
	"api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SupportTicketController struct {
	supportTicketService *services.SupportTicketService
}

func NewSupportTicketController() *SupportTicketController {
	return &SupportTicketController{
		supportTicketService: services.NewSupportTicketService(),
	}
}

func (c *SupportTicketController) ApproveSupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	ticketID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	isAllowed, err := c.supportTicketService.ApproveSupportTicket(ctx, ticketID)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isAllowed {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(403, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *SupportTicketController) DenySupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	ticketID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	isAllowed, err := c.supportTicketService.DenySupportTicket(ctx, ticketID)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isAllowed {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(403, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(200, response)
}
