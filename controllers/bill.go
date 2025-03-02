package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"

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

func (c *BillController) GetBill(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var bill = &[]models.BillModel{}

	isAuthenticated, err := c.billService.GetBill(ctx, bill)
	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isAuthenticated {
		response.Message = config.GetMessageCode("INVALID_CREDENTIALS")
		ctx.JSON(401, response)
		return
	}

	response.Data = bill
	ctx.JSON(200, response)
}


func (c *BillController) CreateBill(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	var bill = &structs.NewBill{}

	if err := ctx.ShouldBind(bill); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(bill); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	// if err := c.buildingService.CreateBuilding(ctx, building); err != nil {
	// 	response.Message = config.GetMessageCode("SYSTEM_ERROR")
	// 	ctx.JSON(500, response)
	// 	return
	// }

	ctx.JSON(200, response)
}
