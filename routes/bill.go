package routes

import (
	"api/controllers"
	"api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBillRoutes(router *gin.RouterGroup) {
	r := router.Group("/bill")
	billController := controllers.NewBillController()
	authorizationMiddle := middlewares.NewAuthorizationMiddleware()

	r.Use(authorizationMiddle.AuthManagerMiddleware)
	{
		r.GET("/", billController.GetBill)
	}

	r.Use(authorizationMiddle.AuthOwnerMiddleware)
	{
		r.POST("/add", billController.CreateBill)
	}
}
