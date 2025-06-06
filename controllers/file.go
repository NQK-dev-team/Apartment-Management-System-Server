package controllers

import (
	"api/config"
	"api/services"
	"api/utils"
	"net/http"

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

	buildingID := ctx.Param("buildingID")
	filename := ctx.Param("fileName")

	if buildingID == "" || filename == "" {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	// file := &structs.CustomFileStruct{}
	// if err := utils.GetFile(file, "images/buildings/"+buildingID+"/"+filename); err != nil {
	// 	response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
	// 	ctx.JSON(http.StatusNotFound, response)
	// 	return
	// }

	// ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	// ctx.Header("Content-Type", "image/"+ext)
	// ctx.Header("Content-Disposition", "inline; filename="+file.Filename)
	// ctx.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// if _, err := io.Copy(ctx.Writer, bytes.NewReader(file.Content)); err != nil {
	// 	response.Message = config.GetMessageCode("SYSTEM_ERROR")
	// 	ctx.JSON(http.StatusInternalServerError, response)
	// 	return
	// }

	if err := utils.GetFile(ctx, "images/buildings/"+buildingID+"/"+filename); err != nil {
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

	// file := &structs.CustomFileStruct{}
	// fmt.Println("images/buildings/" + buildingID + "/rooms/" + roomNo + "/" + filename)
	// if err := utils.GetFile(file, "images/buildings/"+buildingID+"/rooms/"+roomNo+"/"+filename); err != nil {
	// 	response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
	// 	ctx.JSON(http.StatusNotFound, response)
	// 	return
	// }

	// ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	// ctx.Header("Content-Type", "image/"+ext)
	// ctx.Header("Content-Disposition", "inline; filename="+file.Filename)
	// ctx.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// if _, err := io.Copy(ctx.Writer, bytes.NewReader(file.Content)); err != nil {
	// 	response.Message = config.GetMessageCode("SYSTEM_ERROR")
	// 	ctx.JSON(http.StatusInternalServerError, response)
	// 	return
	// }

	if err := utils.GetFile(ctx, "images/buildings/"+buildingID+"/rooms/"+roomNo+"/"+filename); err != nil {
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

	// file := &structs.CustomFileStruct{}
	// fmt.Println("images/buildings/" + buildingID + "/rooms/" + roomNo + "/" + filename)
	// if err := utils.GetFile(file, "images/buildings/"+buildingID+"/rooms/"+roomNo+"/"+filename); err != nil {
	// 	response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
	// 	ctx.JSON(http.StatusNotFound, response)
	// 	return
	// }

	// ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	// ctx.Header("Content-Type", "image/"+ext)
	// ctx.Header("Content-Disposition", "inline; filename="+file.Filename)
	// ctx.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// if _, err := io.Copy(ctx.Writer, bytes.NewReader(file.Content)); err != nil {
	// 	response.Message = config.GetMessageCode("SYSTEM_ERROR")
	// 	ctx.JSON(http.StatusInternalServerError, response)
	// 	return
	// }

	if err := utils.GetFile(ctx, "images/users/"+userID+"/"+filename); err != nil {
		response.Message = config.GetMessageCode("IMAGE_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}
}
