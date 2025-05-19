package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func InitContractRoutes(router *gin.RouterGroup) {
	contractRoutes := router.Group("/contract")
	contractController := controllers.NewContractController()

	contractRoutes.GET("/", contractController.GetContractList)
}
