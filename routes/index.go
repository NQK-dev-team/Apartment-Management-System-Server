package routes

import (
	"api/controllers"
	middlewares "api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.RouterGroup) {
	r := router.Group("v1")

	// Init authentication sub-routes
	InitAuthRoutes(r)

	// Init routes for MoMo payment confirmation without authentication
	InitMoMoBillRoutes(r)

	// Apply the jwtMiddleware to other sub-routes
	authMiddleware := middlewares.NewAuthenticationMiddleware()
	r.Use(authMiddleware.AuthMiddleware)
	{
		// Init other sub-routes
		InitFileRoutes(r)
		InitBuildingRoutes(r)
		InitUserRoutes(r)
		InitSupportTicketRoutes(r)
		InitContractRoutes(r)
		InitRoomRoutes(r)
		InitBillRoutes(r)
		InitNotificationRoutes(r)
		InitUploadRoutes(r)
	}
}

func InitWebSocketRoutes(router *gin.RouterGroup) {
	r := router.Group("v1")
	r.GET("/", controllers.HandleConnection)

	go controllers.HandleBroadcast()
}
