package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoomRoutes(router *gin.RouterGroup) {
	roomRoutes := router.Group("/room")
	roomController := controllers.NewRoomController()
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()

	roomRoutes.Use(authorizationMiddle.AuthCustomerMiddleware)
	{
		roomRoutes.GET("/", roomController.GetRoomList)
		roomRoutes.GET("/:id", roomController.GetRoomByID)
		roomRoutes.GET("/:id/support-ticket", roomController.GetRoomSupportTicket)
	}
}
