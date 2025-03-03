package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBuildingRoutes(router *gin.RouterGroup) {
	r := router.Group("/building")
	buildingController := controllers.NewBuildingController()
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()

	r.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		r.GET("/", buildingController.GetBuilding)
		r.GET("/:id", buildingController.GetBuildingDetail)
		r.GET("/:id/room", buildingController.GetBuildingRoom)
	}

	r.Use(authorizationMiddle.AuthOwnerMiddleware)
	{
		r.POST("/add", buildingController.CreateBuilding)
	}
}
