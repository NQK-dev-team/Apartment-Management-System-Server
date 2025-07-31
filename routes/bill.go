package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBillRoutes(router *gin.RouterGroup) {
	billRoutes := router.Group("/bill")
	billController := controllers.NewBillController()
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()

	billRoutes.GET("/", billController.GetBillList)
	billRoutes.GET("/:id", billController.GetBillDetail)

	billRoutes.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		billRoutes.POST("/delete-many", billController.DeleteManyBills)
		billRoutes.POST("/:id/update", billController.UpdateBill)
		billRoutes.POST("/add", billController.AddBill)
	}
}
