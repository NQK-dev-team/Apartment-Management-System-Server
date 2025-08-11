package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitSupportTicketRoutes(router *gin.RouterGroup) {
	managerRoutes := router.Group("/support-ticket")
	customerRoutes := router.Group("/support-ticket")
	generalRoutes := router.Group("/support-ticket")

	authorizationMiddle := middlewares.NewAuthorizationMiddleware()
	controller := controllers.NewSupportTicketController()

	generalRoutes.GET("/", controller.GetSupportTickets)

	customerRoutes.Use(authorizationMiddle.AuthCustomerMiddleware)
	{
		customerRoutes.POST("/delete-many", controller.DeleteManySupportTickets)
		customerRoutes.POST("/:id/update", controller.UpdateSupportTicket)
	}

	managerRoutes.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		managerRoutes.POST("/:id/approve", controller.ApproveSupportTicket)
		managerRoutes.POST("/:id/deny", controller.DenySupportTicket)
	}
}
