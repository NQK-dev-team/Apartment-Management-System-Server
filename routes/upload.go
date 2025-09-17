package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUploadRoutes(router *gin.RouterGroup) {
	r := router.Group("/upload")
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()
	controller := controllers.NewUploadController()

	r.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		r.POST("/add", controller.UploadFile)
		r.GET("/not-process", controller.GetNotProcessedFiles)
		r.GET("/processed", controller.GetProcessedFiles)
	}
}
