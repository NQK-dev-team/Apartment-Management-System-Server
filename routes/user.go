package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(router *gin.RouterGroup) {
	staffRoutes := router.Group("/staff")
	customerRoutes := router.Group("/customer")
	userRoutes := router.Group("/user")
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
		staffRoutes.POST("/add", userController.AddStaff)
		staffRoutes.POST("/:id/update", userController.UpdateStaff)
	}

	customerRoutes.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		customerRoutes.GET("/", userController.GetCustomerList)
		customerRoutes.POST("/delete-many", userController.DeleteCustomers)
		customerRoutes.GET("/:id", userController.GetCustomerDetail)
		customerRoutes.GET("/:id/contract", userController.GetCustomerContract)
		customerRoutes.GET("/:id/ticket", userController.GetCustomerTicket)
		customerRoutes.POST("/add", userController.AddCustomer)
	}

	userRoutes.GET("/profile", userController.GetUserInfo)
	userRoutes.POST("/profile/update", userController.UpdateUserInfo)
	userRoutes.POST("/security/change-password", userController.ChangePassword)
	userRoutes.POST("/security/change-email", userController.ChangeEmail)
}
