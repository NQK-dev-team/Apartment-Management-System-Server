package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitFileRoutes(router *gin.RouterGroup) {
	imageRoutes := router.Group("/images")
	// fileRoutes := router.Group("/files")

	fileController := controllers.NewFileController()

	imageRoutes.GET("/buildings/:buildingID/rooms/:roomNo/:fileName", fileController.GetRoomImage)
	imageRoutes.GET("/buildings/:buildingID/:fileName", fileController.GetBuildingImage)
	imageRoutes.GET("/users/:userID/:fileName", fileController.GetUserImage)
}
