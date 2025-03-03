package controllers

import (
	"api/config"
	"api/services"
	"api/structs"
	"api/utils"
	"bytes"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	authenticationService *services.AuthenticationService
}

func NewFileController() *FileController {
	return &FileController{
		authenticationService: services.NewAuthenticationService(),
	}
}

func (c *FileController) GetBuildingImage(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id := ctx.Param("id")
	filename := ctx.Param("filename")

	if id == "" || filename == "" {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	file := &structs.CustomFileStruct{}

	if err := utils.GetFile(file, "images/buildings/"+id+"/"+filename); err != nil {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	ctx.Header("Content-Type", "image/"+ext)
	ctx.Header("Content-Disposition", "inline; filename="+file.Filename)
	ctx.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	if _, err := io.Copy(ctx.Writer, bytes.NewReader(file.Content)); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}
}
