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
}

func NewFileController() *FileController {
	return &FileController{
		authenticationService: services.NewAuthenticationService(),
		contractService:       services.NewContractService(),
		supportTicketService:  services.NewSupportTicketService(),
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
