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

type FileController struct {
	authenticationService *services.AuthenticationService
	contractService       *services.ContractService
	supportTicketService  *services.SupportTicketService
	notificationService   *services.NotificationService
	uploadService         *services.UploadService
}

func NewFileController() *FileController {
	return &FileController{
		authenticationService: services.NewAuthenticationService(),
		contractService:       services.NewContractService(),
		supportTicketService:  services.NewSupportTicketService(),
		notificationService:   services.NewNotificationService(),
		uploadService:         services.NewUploadService(),
	}
}

func (c *FileController) GetBuildingImage(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID := ctx.Param("buildingID")
	filename := ctx.Param("fileName")

	if buildingID == "" || filename == "" {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	if err := utils.GetFile(ctx, constants.GetBuildingImageURL("images", buildingID, filename)); err != nil {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}

func (c *FileController) GetRoomImage(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	buildingID := ctx.Param("buildingID")
	roomNo := ctx.Param("roomNo")
	filename := ctx.Param("fileName")

	if buildingID == "" || roomNo == "" || filename == "" {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	if err := utils.GetFile(ctx, constants.GetRoomImageURL("images", buildingID, roomNo, filename)); err != nil {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}

func (c *FileController) GetUserImage(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	userID := ctx.Param("userID")
	filename := ctx.Param("fileName")

	if userID == "" || filename == "" {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	role := ctx.GetString("role")
	if role == constants.Roles.Customer && userID != strconv.FormatInt(ctx.GetInt64("userID"), 10) {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	if err := utils.GetFile(ctx, constants.GetUserImageURL("images", userID, filename)); err != nil {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}

func (c *FileController) GetContractFile(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	contractIDStr := ctx.Param("contractID")
	contractID, _ := strconv.ParseInt(contractIDStr, 10, 64)
	filename := ctx.Param("fileName")

	if contractID == 0 || filename == "" {
		response.Message = config.GetMessageCode("FILE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	contract := &structs.Contract{}

	isAllowed, err := c.contractService.GetContractDetail(ctx, contract, contractID)
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

	if err := utils.GetFile(ctx, constants.GetContractFileURL("files", contractIDStr, filename)); err != nil {
		response.Message = config.GetMessageCode("FILE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}

func (c *FileController) GetTicketImage(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	ticketIDStr := ctx.Param("ticketID")
	ticketID, _ := strconv.ParseInt(ticketIDStr, 10, 64)
	filename := ctx.Param("fileName")

	if ticketID == 0 || filename == "" {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ticket := &structs.SupportTicket{}
	if err := c.supportTicketService.GetSupportTicket(ctx, ticket, ticketID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if ticket.ID == 0 {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	if err := utils.GetFile(ctx, constants.GetTicketImageURL("images", ticketIDStr, filename)); err != nil {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}

func (c *FileController) GetNotificationFile(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	notificationIDStr := ctx.Param("notificationID")
	notificationID, _ := strconv.ParseInt(notificationIDStr, 10, 64)
	filename := ctx.Param("fileName")

	if notificationID == 0 || filename == "" {
		response.Message = config.GetMessageCode("FILE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	isAllowed, err := c.notificationService.CheckUserGetNotification(ctx, notificationID)
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

	if err := utils.GetFile(ctx, constants.GetNotificationFileURL("files", notificationIDStr, filename)); err != nil {
		response.Message = config.GetMessageCode("FILE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}

func (c *FileController) GetUploadFile(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	role := ctx.GetString("role")

	if role != constants.Roles.Manager && role != constants.Roles.Owner {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	uploadIDStr := ctx.Param("uploadID")
	filename := ctx.Param("fileName")

	if err := utils.GetFile(ctx, constants.GetUploadFileURL("files", uploadIDStr, filename)); err != nil {
		response.Message = config.GetMessageCode("FILE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}
