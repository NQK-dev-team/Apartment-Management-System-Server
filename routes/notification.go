package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitNotificationRoutes(router *gin.RouterGroup) {
	r := router.Group("/notification")
	controller := controllers.NewNotificationController()
	authorizationMiddleware := middlewares.NewAuthorizationMiddleware()

	r.GET("/inbox", controller.GetInbox)
	r.GET("/marked", controller.GetMarked)
	r.PATCH("/:id/read", controller.MarkAsRead)
	r.PATCH("/:id/unread", controller.MarkAsUnread)
	r.PATCH("/:id/mark", controller.MarkAsImportant)
	r.PATCH("/:id/unmark", controller.UnmarkAsImportant)

	r.Use(authorizationMiddleware.AuthManagerMiddleware)
	{
		r.POST("/add", controller.AddNotification)
		r.GET("/sent", controller.GetSent)
		// r.DELETE("/:id", controller.DeleteNotification)
		// r.GET("/:id", controller.GetNotificationDetail)
	}
}
