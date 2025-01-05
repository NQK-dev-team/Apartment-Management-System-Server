package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(router *gin.RouterGroup) {
	r := router.Group("/authentication")
	authenticationController := controllers.NewAuthenticationController()
	r.POST("/login", authenticationController.Login)
	r.POST("/logout", authenticationController.Logout)
	r.POST("/recovery", authenticationController.Recovery)
	r.POST("/verify-token", authenticationController.VerifyToken)
	r.POST("/refresh-token", authenticationController.RefreshToken)
}
