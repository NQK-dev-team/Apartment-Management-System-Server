package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(router *gin.RouterGroup) {
	staffRoutes := router.Group("/staff")
	// customerRoutes := router.Group("/customer")
	userController := controllers.NewUserController()

	// authorizationMiddle := middlewares.NewAuthorizationMiddleware()

	staffRoutes.GET("/", userController.GetStaffList)
}
