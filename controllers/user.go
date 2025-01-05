package controllers

import (
	"api/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{userService: services.NewUserService()}
}

func (c *UserController) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Get",
	})
}
