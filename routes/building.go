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
		r.GET("/:id/", buildingController.GetBuildingDetail)
		r.DELETE("/:id", buildingController.DeleteBuilding)
		r.GET("/:id/room", buildingController.GetBuildingRoom)
		r.GET("/:id/service", buildingController.GetBuildingService)
		r.POST("/:id/deleteRooms", buildingController.DeleteRooms)
		r.POST("/:id/deleteServices", buildingController.DeleteServices)
		r.POST("/:id/service/add", buildingController.AddService)
		r.POST("/:id/service/:serviceID/edit", buildingController.EditService)
		r.POST("/:id/room/add", buildingController.AddRoom)
		r.GET("/:id/schedule", buildingController.GetBuildingSchedule)
		r.POST("/:id/update", buildingController.UpdateBuilding)
	}

	r.Use(authorizationMiddle.AuthOwnerMiddleware)
	{
		r.POST("/add", buildingController.CreateBuilding)
	}
}
