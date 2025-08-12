package controllers

import (
	"api/config"
	"api/constants"
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
		startDate = utils.GetFirstDayOfMonth("")
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

func (c *SupportTicketController) DeleteManySupportTickets(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	type deleteIDs struct {
		IDs []int64 `json:"IDs" validate:"required"`
	}

	input := &deleteIDs{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := constants.Validate.Struct(input); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	validID, err := c.supportTicketService.DeleteTickets(ctx, input.IDs)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !validID {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *SupportTicketController) UpdateSupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	ticketID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ticketID = 0
	}

	ticket := &structs.UpdateSupportTicketRequest{}
	if err := ctx.ShouldBind(ticket); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	form, _ := ctx.MultipartForm()
	ticket.NewFiles = form.File["newFiles[]"]

	if err := constants.Validate.Struct(ticket); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	fileValidation := &structs.ValidateSupportTicketFile{
		Images: []structs.ImageValidation{},
	}

	for _, file := range ticket.NewFiles {
		fileValidation.Images = append(fileValidation.Images, structs.ImageValidation{
			Type: file.Header.Get("Content-Type"),
			Size: file.Size,
		})
	}

	if err := constants.Validate.Struct(fileValidation); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, isFound, err := c.supportTicketService.UpdateSupportTicket(ctx, ticketID, ticket)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isFound {
		response.Message = config.GetMessageCode("DATA_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
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

func (c *SupportTicketController) AddSupportTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	ticket := &structs.CreateSupportTicketRequest{}
	if err := ctx.ShouldBind(ticket); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	form, _ := ctx.MultipartForm()
	ticket.Files = form.File["files[]"]

	if err := constants.Validate.Struct(ticket); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	fileValidation := &structs.ValidateSupportTicketFile{
		Images: []structs.ImageValidation{},
	}

	for _, file := range ticket.Files {
		fileValidation.Images = append(fileValidation.Images, structs.ImageValidation{
			Type: file.Header.Get("Content-Type"),
			Size: file.Size,
		})
	}

	if err := constants.Validate.Struct(fileValidation); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, err := c.supportTicketService.AddSupportTicket(ctx, ticket)
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

	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
