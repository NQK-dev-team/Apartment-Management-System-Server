package controllers

import (
	"api/config"
	"api/models"
	"api/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{userService: services.NewUserService()}
}

func (c *UserController) GetStaffList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	users := &[]models.UserModel{}

	if err := c.userService.GetStaffList(ctx, users); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = users
	ctx.JSON(200, response)
}
