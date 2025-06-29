package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContractController struct {
	contractService *services.ContractService
}

func NewContractController() *ContractController {
	return &ContractController{
		contractService: services.NewContractService(),
	}
}

func (c *ContractController) GetContractList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	contracts := []structs.Contract{}

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

	if err := c.contractService.GetContractList(ctx, &contracts, limit, offset); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = contracts
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *ContractController) GetContractDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	contract := &structs.Contract{}

	isAllowed, err := c.contractService.GetContractDetail(ctx, contract, id)

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

	response.Data = contract
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *ContractController) GetContractBill(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	bills := []models.BillModel{}

	isAllowed, err := c.contractService.GetContractBill(ctx, &bills, id)

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

	response.Data = bills
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *ContractController) DeleteManyContracts(ctx *gin.Context) {
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

	validID, err := c.contractService.DeleteContract2(ctx, input.IDs)
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

func (c *ContractController) UpdateContract(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		id = 0
	}

	contract := &structs.EditContract{}

	if err := ctx.ShouldBind(contract); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	for idx := range contract.TotalNewfiles {
		fileStr := fmt.Sprintf("file[%d]file", idx)
		titleStr := fmt.Sprintf("file[%d]title", idx)

		fileHeader, _ := ctx.FormFile(fileStr)
		title := ctx.PostForm(titleStr)

		contract.NewFiles = append(contract.NewFiles, structs.ContractFile{
			File:  fileHeader,
			Title: title,
		})
	}

	if err := constants.Validate.Struct(contract); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	validateContractFiles := &structs.ValidateEditContractFile{
		NewFiles: []structs.FileValidation{},
	}

	for _, file := range contract.NewFiles {
		validateContractFiles.NewFiles = append(validateContractFiles.NewFiles, structs.FileValidation{
			Type: file.File.Header.Get("Content-Type"),
			Size: file.File.Size,
		})
	}

	if err := constants.Validate.Struct(validateContractFiles); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, isValidUpdate, err := c.contractService.UpdateContract(ctx, contract, id)

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

	if !isValidUpdate {
		response.Message = config.GetMessageCode("UPDATE_FAILED")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
