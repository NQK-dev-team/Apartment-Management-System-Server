package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(router *gin.RouterGroup) {
	staffRoutes := router.Group("/staff")
	// customerRoutes := router.Group("/customer")
	userController := controllers.NewUserController()

	authorizationMiddle := middlewares.NewAuthorizationMiddleware()
	staffRoutes.Use(authorizationMiddle.AuthOwnerMiddleware)
	{
		staffRoutes.GET("/", userController.GetStaffList)
	}
}
