package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		notificationService: services.NewNotificationService(),
	}
}

func (c *NotificationController) GetInbox(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	notifications := &[]structs.Notification{}

	if err := c.notificationService.GetInboxNotifications(ctx, notifications, limit, offset); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = notifications
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) GetMarked(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	notifications := &[]structs.Notification{}

	if err := c.notificationService.GetMarkedNotifications(ctx, notifications, limit, offset); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = notifications
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) GetSent(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	notifications := &[]models.NotificationModel{}

	if err := c.notificationService.GetSentNotifications(ctx, notifications, limit, offset); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = notifications
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) AddNotification(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	newNotification := &structs.NewNotification{}

	if err := ctx.ShouldBind(newNotification); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	newNotification.Title = strings.TrimSpace(newNotification.Title)
	newNotification.Content = strings.TrimSpace(newNotification.Content)

	// Replace any whitespace in the receiver string
	newNotification.ReceiverStr = strings.ReplaceAll(newNotification.ReceiverStr, " ", "")
	// Convert ReceiverStr data to Receivers of []int64
	receiverIDs := strings.Split(newNotification.ReceiverStr, ",")
	for _, idStr := range receiverIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			newNotification.Receivers = append(newNotification.Receivers, id)
		}
	}

	form, _ := ctx.MultipartForm()
	newNotification.Files = form.File["files[]"]

	if err := constants.Validate.Struct(newNotification); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	fileValidation := &structs.ValidateNoticeFile{
		Files: []structs.FileValidation{},
	}

	for _, file := range newNotification.Files {
		fileValidation.Files = append(fileValidation.Files, structs.FileValidation{
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

	if err := c.notificationService.AddNotification(ctx, newNotification); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	signal := &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewInbox,
		},
		Users: newNotification.Receivers,
	}

	AddBroadcast(signal)

	signal = &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewSent,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		id = 0
	}

	isAllowed, err := c.notificationService.UpdateNotificationReadStatus(ctx, id, true)
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

	signal := &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewImportant,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	signal = &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewInbox,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) MarkAsUnread(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		id = 0
	}

	isAllowed, err := c.notificationService.UpdateNotificationReadStatus(ctx, id, false)
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

	signal := &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewImportant,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	signal = &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewInbox,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) MarkAsImportant(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		id = 0
	}

	isAllowed, err := c.notificationService.UpdateNotificationImportantStatus(ctx, id, true)
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

	signal := &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewImportant,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	signal = &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewInbox,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) UnmarkAsImportant(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		id = 0
	}

	isAllowed, err := c.notificationService.UpdateNotificationImportantStatus(ctx, id, false)
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

	signal := &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewImportant,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	signal = &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewInbox,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *NotificationController) MarkMultiNotiAsRead(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	idsStruct := &structs.IDList{}
	if err := ctx.ShouldBindJSON(idsStruct); err != nil {
		response.Message = config.GetMessageCode("INVALID_REQUEST")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := constants.Validate.Struct(idsStruct); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, err := c.notificationService.MarkMultiNotiAsRead(ctx, &idsStruct.IDs)
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

	signal := &structs.NotificationWS{
		BaseWSStruct: structs.BaseWSStruct{
			Type: constants.Common.WebsocketSignalType.NewInbox,
		},
		Users: []int64{ctx.GetInt64("userID")},
	}

	AddBroadcast(signal)

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

// func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
// 	response := config.NewDataResponse(ctx)

// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)

// 	if err != nil {
// 		id = 0
// 	}

// 	receivers := []int64{}

// 	isAllowed, err := c.notificationService.DeleteNotification(ctx, id, &receivers)
// 	if err != nil {
// 		response.Message = config.GetMessageCode("SYSTEM_ERROR")
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	if !isAllowed {
// 		response.Message = config.GetMessageCode("PERMISSION_DENIED")
// 		ctx.JSON(http.StatusForbidden, response)
// 		return
// 	}

// 	signal := &structs.NotificationWS{
// 		BaseWSStruct: structs.BaseWSStruct{
// 			Type: constants.Common.WebsocketSignalType.NewImportant,
// 		},
// 		Users: receivers,
// 	}

// 	AddBroadcast(signal)

// 	signal = &structs.NotificationWS{
// 		BaseWSStruct: structs.BaseWSStruct{
// 			Type: constants.Common.WebsocketSignalType.NewInbox,
// 		},
// 		Users: receivers,
// 	}

// 	AddBroadcast(signal)

// 	response.Message = config.GetMessageCode("DELETE_SUCCESS")
// 	ctx.JSON(http.StatusOK, response)
// }

// func (c *NotificationController) GetNotificationDetail(ctx *gin.Context) {
// 	response := config.NewDataResponse(ctx)

// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)

// 	if err != nil {
// 		id = 0
// 	}

// 	notification := &structs.NotificationDetail{}

// 	isFound, isAllowed, err := c.notificationService.GetNotificationDetail(ctx, id, notification)
// 	if err != nil {
// 		response.Message = config.GetMessageCode("SYSTEM_ERROR")
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	if !isAllowed {
// 		response.Message = config.GetMessageCode("PERMISSION_DENIED")
// 		ctx.JSON(http.StatusForbidden, response)
// 		return
// 	}

// 	if !isFound {
// 		response.Message = config.GetMessageCode("DATA_NOT_FOUND")
// 		ctx.JSON(http.StatusNotFound, response)
// 		return
// 	}

// 	response.Data = notification
// 	response.Message = config.GetMessageCode("GET_SUCCESS")
// 	ctx.JSON(http.StatusOK, response)
// }
