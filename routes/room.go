package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoomRoutes(router *gin.RouterGroup) {
	roomRoutes := router.Group("/room")
	roomController := controllers.NewRoomController()

	roomRoutes.GET("/", roomController.GetRoomList)
}
