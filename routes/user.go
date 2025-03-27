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

	staffRoutes.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		staffRoutes.GET("/", userController.GetStaffList)
	}

	staffRoutes.Use(authorizationMiddle.AuthOwnerMiddleware)
	{
		staffRoutes.GET("/:id", userController.GetStaffDetail)
		staffRoutes.GET("/:id/schedule", userController.GetStaffSchedule)
		staffRoutes.GET("/:id/contract", userController.GetStaffRelatedContract)
		staffRoutes.GET("/:id/ticket", userController.GetStaffRelatedTicket)
		staffRoutes.POST("/delete-many", userController.DeleteStaffs)
	}
}
