package controllers

import (
	"api/config"
	"api/services"
	"api/structs"
	"api/utils"
	"net/http"
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

func (c *SupportTicketController) GetSupportTickets(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	if startDate == "" {
		startDate = utils.GetFirstDayOfMonth()
	}

	if endDate == "" {
		endDate = utils.GetCurrentDate()
	}

	tickets := []structs.SupportTicket{}

	if err := c.supportTicketService.GetSupportTickets(ctx, &tickets, limit, offset, startDate, endDate); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = tickets
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *SupportTicketController) ApproveSupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	ticketID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, err := c.supportTicketService.ApproveSupportTicket(ctx, ticketID)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isAllowed {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *SupportTicketController) DenySupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	ticketID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, err := c.supportTicketService.DenySupportTicket(ctx, ticketID)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isAllowed {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
