package controllers

import (
	"api/config"
	"api/constants"
	"api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BillController struct {
	billService *services.BillService
}

func NewBillController() *BillController {
	return &BillController{
		billService: services.NewBillService(),
	}
}

func (c *BillController) DeleteManyBills(ctx *gin.Context) {
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

	isValid, err := c.billService.DeleteBill(ctx, input.IDs)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
