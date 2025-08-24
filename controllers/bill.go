package controllers

import (
	"api/config"
	"api/constants"
	"api/services"
	"api/structs"
	"api/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (c *BillController) GetBillList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")
	startMonth := ctx.Query("startMonth")
	endMonth := ctx.Query("endMonth")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	if startMonth == "" {
		startMonth = utils.GetFirstDayOfQuarter()
	} else {
		if _, err := time.Parse("2006-01", startMonth); err != nil {
			startMonth = utils.GetFirstDayOfQuarter()
		} else {
			startMonth = utils.GetFirstDayOfMonth(startMonth)
		}
	}

	if endMonth != "" {
		if _, err := time.Parse("2006-01", endMonth); err != nil {
			endMonth = utils.GetLastDayOfMonth("")
		} else {
			endMonth = utils.GetLastDayOfMonth(endMonth)
		}
	} else {
		endMonth = utils.GetLastDayOfMonth(endMonth)
	}

	bills := []structs.Bill{}

	if err := c.billService.GetBillList(ctx, &bills, limit, offset, startMonth, endMonth); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = bills
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
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

func (c *BillController) GetBillDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	bill := &structs.Bill{}

	isAllowed, err := c.billService.GetBillDetail(ctx, bill, id)

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

	response.Data = bill
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *BillController) UpdateBill(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	bill := structs.UpdateBill{}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		id = 0
	}

	if err := ctx.ShouldBindJSON(&bill); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	bill.Title = strings.TrimSpace(bill.Title)
	bill.Note = strings.TrimSpace(bill.Note)
	for index, _ := range bill.Payments {
		bill.Payments[index].Note = strings.TrimSpace(bill.Payments[index].Note)
		bill.Payments[index].Name = strings.TrimSpace(bill.Payments[index].Name)
	}
	for index, _ := range bill.NewPayments {
		bill.NewPayments[index].Note = strings.TrimSpace(bill.NewPayments[index].Note)
		bill.NewPayments[index].Name = strings.TrimSpace(bill.NewPayments[index].Name)
	}

	if err := constants.Validate.Struct(bill); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isAllowed, isValid, err := c.billService.UpdateBill(ctx, &bill, id)
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

	if !isValid {
		response.Message = config.GetMessageCode("UPDATE_FAILED")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *BillController) AddBill(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	bill := structs.AddBill{}

	if err := ctx.ShouldBindJSON(&bill); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	bill.Title = strings.TrimSpace(bill.Title)
	bill.Note = strings.TrimSpace(bill.Note)
	for index, _ := range bill.BillPayments {
		bill.BillPayments[index].Note = strings.TrimSpace(bill.BillPayments[index].Note)
		bill.BillPayments[index].Name = strings.TrimSpace(bill.BillPayments[index].Name)
	}

	if err := constants.Validate.Struct(bill); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var newBillID int64

	isAllowed, isValid, err := c.billService.AddBill(ctx, &bill, &newBillID)
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

	if !isValid {
		response.Message = config.GetMessageCode("CREATE_FAILED")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Data = newBillID
	response.Message = config.GetMessageCode("ADD_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
