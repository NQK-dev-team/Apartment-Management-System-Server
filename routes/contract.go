package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitContractRoutes(router *gin.RouterGroup) {
	contractRoutes := router.Group("/contract")
	contractController := controllers.NewContractController()
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()

	contractRoutes.GET("/", contractController.GetContractList)
	contractRoutes.GET("/:id", contractController.GetContractDetail)
	contractRoutes.GET("/:id/bill", contractController.GetContractBill)

	contractRoutes.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		contractRoutes.POST("/delete-many", contractController.DeleteManyContracts)
		contractRoutes.POST("/:id/update", contractController.UpdateContract)
		contractRoutes.POST("/add", contractController.AddContract)
	}
}
