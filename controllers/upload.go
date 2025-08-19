package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	service *services.UploadService
}

func NewUploadController() *UploadController {
	return &UploadController{
		service: services.NewUploadService(),
	}
}

func (c *UploadController) UploadFile(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	upload := &structs.UploadStruct{}

	if err := ctx.ShouldBind(upload); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := constants.Validate.Struct(upload); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	fileValidation := &structs.ValidateUploadFile{
		File: structs.UploadValidation{
			Type: upload.File.Header.Get("Content-Type"),
			Size: upload.File.Size,
		},
	}

	if err := constants.Validate.Struct(fileValidation); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := c.service.UploadFile(ctx, upload); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *UploadController) GetNotProcessedFiles(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	uploads := &[]models.UploadFileModel{}

	typeStr := ctx.DefaultQuery("type", "0")

	uploadType, err := strconv.Atoi(typeStr)

	if err != nil {
		uploadType = 0
	}

	if err := c.service.GetUploads(ctx, uploads, uploadType, false); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = uploads
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *UploadController) GetProcessedFiles(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	uploads := &[]models.UploadFileModel{}

	typeStr := ctx.DefaultQuery("type", "0")

	uploadType, err := strconv.Atoi(typeStr)

	if err != nil {
		uploadType = 0
	}

	if err := c.service.GetUploads(ctx, uploads, uploadType, true); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = uploads
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
