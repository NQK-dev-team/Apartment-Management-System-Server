package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitSupportTicketRoutes(router *gin.RouterGroup) {
	r := router.Group("/support-ticket")
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()
	controller := controllers.NewSupportTicketController()

	r.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		r.GET("/", controller.GetSupportTickets)
		r.POST("/:id/approve", controller.ApproveSupportTicket)
		r.POST("/:id/deny", controller.DenySupportTicket)
	}
}
