package controllers

import (
	"api/config"
	"api/services"
	"api/structs"
	"api/utils"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authenticationService *services.AuthenticationService
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{authenticationService: services.NewAuthenticationService()}
}

func (c *AuthenticationController) Login(ctx *gin.Context) {
	account := structs.LoginAccount{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&account); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&account); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		ctx.JSON(400, response)
		return
	}

	jwtToken, refreshToken, err, isEmailVerified := c.authenticationService.Login(ctx, account.Email, account.Password, account.Remember)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if jwtToken == "" && isEmailVerified {
		response.Message = config.GetMessageCode("INVALID_CREDENTIALS")
		ctx.JSON(401, response)
		return
	}

	if jwtToken == "" && !isEmailVerified {
		response.Message = config.GetMessageCode("EMAIL_NOT_VERIFIED")
		ctx.JSON(401, response)
		return
	}

	response.Message = config.GetMessageCode("LOGIN_SUCCESS")
	response.JWTToken = jwtToken
	response.RefreshToken = refreshToken

	ctx.JSON(200, response)
}

func (c *AuthenticationController) VerifyToken(ctx *gin.Context) {
	token := structs.VerifyToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		ctx.JSON(400, response)
		return
	}

	isValid, err := c.authenticationService.VerifyToken(ctx, token.JWTToken)

	if err != nil {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(401, response)
		return
	}

	response.Message = config.GetMessageCode("TOKEN_VERIFY_SUCCESS")
	response.Data = isValid

	ctx.JSON(200, response)
}

func (c *AuthenticationController) RefreshToken(ctx *gin.Context) {
	token := structs.RefreshToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		ctx.JSON(400, response)
		return
	}

	jwtToken, err := c.authenticationService.GetNewToken(ctx, token.RefreshToken)

	if err != nil {
		response.Message = config.GetMessageCode("TOKEN_REFRESH_FAILED")
		ctx.JSON(401, response)
		return
	}

	response.Message = config.GetMessageCode("TOKEN_REFRESH_SUCCESS")
	response.JWTToken = jwtToken

	ctx.JSON(200, response)
}

func (c *AuthenticationController) Logout(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Logout",
	})
}

func (c *AuthenticationController) Recovery(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Recovery",
	})
}
