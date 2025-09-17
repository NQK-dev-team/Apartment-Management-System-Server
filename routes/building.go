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
		r.GET("/:id/statistic", buildingController.GetBuildingStatistic)
		r.GET("/:id/", buildingController.GetBuildingDetail)
		r.GET("/:id/schedule", buildingController.GetBuildingSchedule)
		r.GET("/:id/room", buildingController.GetBuildingRoom)
		r.POST("/:id/update", buildingController.UpdateBuilding)
		r.GET("/:id/room/:roomID", buildingController.GetRoomDetail)
		r.POST("/:id/room/:roomID/update", buildingController.UpdateRoomInformation)
		r.GET("/:id/room/:roomID/contracts", buildingController.GetRoomContract)
		r.GET("/:id/room/:roomID/tickets", buildingController.GetRoomTicket)
		r.POST("/:id/room/:roomID/delete-contracts", buildingController.DeleteRoomContract)
	}

	r.Use(authorizationMiddle.AuthOwnerMiddleware)
	{
		r.POST("/add", buildingController.CreateBuilding)
		r.DELETE("/:id", buildingController.DeleteBuilding)
		r.GET("/statistic", buildingController.GetAllBuildingStatistic)
	}
}
