package routes

import (
	"api/constants"
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitFileRoutes(router *gin.RouterGroup) {
	imageRoutes := router.Group("/images")
	fileRoutes := router.Group("/files")
	fileController := controllers.NewFileController()

	imageRoutes.GET(constants.GetRoomImageURL("", ":buildingID", ":roomNo", ":fileName"), fileController.GetRoomImage)
	imageRoutes.GET(constants.GetBuildingImageURL("", ":buildingID", ":fileName"), fileController.GetBuildingImage)
	imageRoutes.GET(constants.GetUserImageURL("", ":userID", ":fileName"), fileController.GetUserImage)

	fileRoutes.GET(constants.GetContractFileURL("", ":contractID", ":fileName"), fileController.GetContractFile)
}
