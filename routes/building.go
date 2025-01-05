package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitBuildingRoutes(router *gin.RouterGroup) {
	r := router.Group("/building")
	buildingController := controllers.NewBuildingController()

	r.GET("/", buildingController.Get)
}
