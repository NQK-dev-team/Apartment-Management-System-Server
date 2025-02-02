package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authenticationService *services.AuthenticationService
	emailService          *services.EmailService
	userSerivce           *services.UserService
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{
		authenticationService: services.NewAuthenticationService(),
		emailService:          services.NewEmailService(),
		userSerivce:           services.NewUserService(),
	}
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
		response.ValidateError = err.Error()
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
		c.emailService.SendEmailVerificationEmail(ctx, account.Email)
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
		response.ValidateError = err.Error()
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
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	jwtToken, err := c.authenticationService.GetNewToken(ctx, token.RefreshToken)

	if err != nil || jwtToken == "" {
		response.Message = config.GetMessageCode("TOKEN_REFRESH_FAILED")
		ctx.JSON(401, response)
		return
	}

	response.Message = config.GetMessageCode("TOKEN_REFRESH_SUCCESS")
	response.JWTToken = jwtToken

	ctx.JSON(200, response)
}

func (c *AuthenticationController) Recovery(ctx *gin.Context) {
	email := structs.RecoveryEmail{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&email); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&email); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	var user = &models.UserModel{}
	err := c.userSerivce.GetUserByEmail(ctx, email.Email, user)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if user.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	isSpam, err := c.emailService.SendResetPasswordEmail(ctx, email.Email)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if isSpam {
		response.Message = config.GetMessageCode("REQUEST_SPAM")
		ctx.JSON(429, response)
		return
	}

	response.Message = config.GetMessageCode("EMAIL_SENT")
	ctx.JSON(200, response)
}

func (c *AuthenticationController) CheckResetPasswordToken(ctx *gin.Context) {
	token := structs.ResetPasswordToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	isValid, err := c.authenticationService.CheckResetPasswordToken(ctx, token.Token, token.Email)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(401, response)
		return
	}

	response.Message = config.GetMessageCode("TOKEN_VERIFY_SUCCESS")
	ctx.JSON(200, response)
}

func (c *AuthenticationController) ResetPassword(ctx *gin.Context) {
	resetPassword := structs.ResetPassword{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&resetPassword); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&resetPassword); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	isValid, err := c.authenticationService.CheckResetPasswordToken(ctx, resetPassword.Token, resetPassword.Email)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(401, response)
		return
	}
	var user = &models.UserModel{}
	err = c.userSerivce.GetUserByEmail(ctx, resetPassword.Email, user)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if user.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	user.Password = resetPassword.Password
	ctx.Set("userID", user.ID)

	if err := c.userSerivce.UpdateUser(ctx, user); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if err := c.authenticationService.DeletePasswordResetToken(ctx, user.Email); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("PASSWORD_RESET")
	ctx.JSON(200, response)
}

func (c *AuthenticationController) VerifyEmail(ctx *gin.Context) {
	token := structs.VerifyEmailToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := utils.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	var user = &models.UserModel{}
	err := c.userSerivce.GetUserByEmail(ctx, token.Email, user)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if user.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(404, response)
		return
	}

	ctx.Set("userID", user.ID)

	isValid, err := c.authenticationService.VerifyEmail(ctx, token)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(401, response)
		return
	}

	response.Message = config.GetMessageCode("TOKEN_VERIFY_SUCCESS")
	ctx.JSON(200, response)
}

func (c *AuthenticationController) Logout(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	// Get the JWT token from the header
	jwt := ctx.GetHeader("Authorization")
	if jwt != "" {
		// Trim the "Bearer " prefix
		jwt = jwt[7:]
	}
	// Get the refresh token from the header
	refreshToken := ctx.GetHeader("X-Refresh-Token")

	if refreshToken != "" {
		claims := c.authenticationService.ExtractJWTData(ctx, jwt)

		if claims == nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.JSON(500, response)
			return
		}

		err := c.authenticationService.DeleteRefreshToken(ctx, claims.UserID)

		if err != nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.JSON(500, response)
			return
		}
	}

	response.Message = config.GetMessageCode("LOGOUT_SUCCESS")
	ctx.JSON(200, response)
}
