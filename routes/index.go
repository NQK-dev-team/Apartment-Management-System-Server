package routes

import (
	middlewares "api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.RouterGroup) {
	r := router.Group("v1")

	// Init authentication sub-routes
	InitAuthRoutes(r)

	// Init file sub-routes
	InitFileRoutes(r)

	// Apply the jwtMiddleware to other sub-routes
	authMiddleware := middlewares.NewAuthenticationMiddleware()
	r.Use(authMiddleware.AuthMiddleware)
	{
		// Init other sub-routes
		InitBuildingRoutes(r)
		InitUserRoutes(r)
		InitSupportTicketRoutes(r)
		InitContractRoutes(r)
		InitRoomRoutes(r)
	}
}
