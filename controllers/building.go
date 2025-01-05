package controllers

import (
	"api/config"
	"api/services"

	"github.com/gin-gonic/gin"
)

type BuildingController struct {
	buildingService *services.BuildingService
}

func NewBuildingController() *BuildingController {
	return &BuildingController{
		buildingService: services.NewBuildingService(),
	}
}

func (c *BuildingController) Get(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	response.ValidateError = "123"
	ctx.JSON(200, response)
}
