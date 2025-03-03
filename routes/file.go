package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitFileRoutes(router *gin.RouterGroup) {
	imageRoutes := router.Group("/images")
	// fileRoutes := router.Group("/files")

	fileController := controllers.NewFileController()

	imageRoutes.GET("/buildings/:id/:filename", fileController.GetBuildingImage)
}
