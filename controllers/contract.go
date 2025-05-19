package controllers

import (
	"api/config"

	"github.com/gin-gonic/gin"
)

type ContractController struct {
}

func NewContractController() *ContractController {
	return &ContractController{}
}

func (c *ContractController) GetContractList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}
