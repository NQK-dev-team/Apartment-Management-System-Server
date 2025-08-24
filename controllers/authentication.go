package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"net/http"
	"strings"

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
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	account.Email = strings.TrimSpace(account.Email)

	if err := constants.Validate.Struct(&account); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	jwtToken, refreshToken, err, isEmailVerified := c.authenticationService.Login(ctx, account.Email, account.Password, account.Remember)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if jwtToken == "" && isEmailVerified {
		response.Message = config.GetMessageCode("INVALID_CREDENTIALS")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	if jwtToken == "" && !isEmailVerified {
		user := &models.UserModel{}

		if err := c.authenticationService.GetUserDataByEmail(ctx, user, account.Email); err != nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		if user.ID == 0 {
			response.Message = config.GetMessageCode("INVALID_CREDENTIALS")
			ctx.JSON(http.StatusUnauthorized, response)
			return
		}

		if !user.VerifiedAfterCreated {
			c.emailService.SendAccountCreatedEmailVerificationEmail(ctx, account.Email)
		} else {
			c.emailService.SendAccountChangeEmailVerificationEmail(ctx, account.Email)
		}
		response.Message = config.GetMessageCode("EMAIL_NOT_VERIFIED")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	response.Message = config.GetMessageCode("LOGIN_SUCCESS")
	response.JWTToken = jwtToken
	response.RefreshToken = refreshToken

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) VerifyToken(ctx *gin.Context) {
	token := structs.VerifyToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := constants.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isValid, err := c.authenticationService.VerifyToken(ctx, token.JWTToken)

	if err != nil {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	response.Message = config.GetMessageCode("VERIFY_SUCCESS")
	response.Data = isValid

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) RefreshToken(ctx *gin.Context) {
	token := structs.RefreshToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := constants.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	jwtToken, err := c.authenticationService.GetNewToken(ctx, token.RefreshToken)

	if err != nil || jwtToken == "" {
		response.Message = config.GetMessageCode("TOKEN_REFRESH_FAILED")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	response.Message = config.GetMessageCode("VERIFY_SUCCESS")
	response.JWTToken = jwtToken

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) Recovery(ctx *gin.Context) {
	email := structs.RecoveryEmail{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&email); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	email.Email = strings.TrimSpace(email.Email)

	if err := constants.Validate.Struct(&email); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var user = &models.UserModel{}
	err := c.userSerivce.GetUserByEmail(ctx, email.Email, user)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if user.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	isSpam, err := c.emailService.SendResetPasswordEmail(ctx, email.Email)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if isSpam {
		response.Message = config.GetMessageCode("REQUEST_SPAM")
		ctx.JSON(http.StatusTooManyRequests, response)
		return
	}

	response.Message = config.GetMessageCode("EMAIL_SENT")
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) CheckResetPasswordToken(ctx *gin.Context) {
	token := structs.ResetPasswordToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token.Email = strings.TrimSpace(token.Email)

	if err := constants.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isValid, err := c.authenticationService.CheckResetPasswordToken(ctx, token.Token, token.Email)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	response.Message = config.GetMessageCode("VERIFY_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) ResetPassword(ctx *gin.Context) {
	resetPassword := structs.ResetPassword{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&resetPassword); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	resetPassword.Email = strings.TrimSpace(resetPassword.Email)

	if err := constants.Validate.Struct(&resetPassword); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isValid, err := c.authenticationService.CheckResetPasswordToken(ctx, resetPassword.Token, resetPassword.Email)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	var user = &models.UserModel{}
	err = c.userSerivce.GetUserByEmail(ctx, resetPassword.Email, user)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if user.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	user.Password = resetPassword.Password
	ctx.Set("userID", user.ID)

	if err := c.userSerivce.UpdateUser(ctx, user); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if err := c.authenticationService.DeletePasswordResetToken(ctx, user.Email); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Message = config.GetMessageCode("PASSWORD_RESET")
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) VerifyEmail(ctx *gin.Context) {
	token := structs.VerifyEmailToken{}
	response := config.NewDataResponse(ctx)

	if err := ctx.ShouldBindJSON(&token); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token.Email = strings.TrimSpace(token.Email)

	if err := constants.Validate.Struct(&token); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isValid, err := c.authenticationService.CheckEmailVerifyToken(ctx, token)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isValid {
		response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	var user = &models.UserModel{}
	err = c.userSerivce.GetUserByEmail(ctx, token.Email, user)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if user.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.Set("userID", user.ID)

	err = c.authenticationService.VerifyEmail(ctx, token)

	if err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Message = config.GetMessageCode("VERIFY_SUCCESS")
	ctx.JSON(http.StatusOK, response)
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
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		err := c.authenticationService.DeleteRefreshToken(ctx, claims.UserID)

		if err != nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	response.Message = config.GetMessageCode("LOGOUT_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthenticationController) VerifyPassword(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	verifyPasswordStruct := structs.VerifyPassword{}

	if err := ctx.ShouldBindJSON(&verifyPasswordStruct); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := constants.Validate.Struct(&verifyPasswordStruct); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = constants.GetValidateErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// Get the JWT token from the header
	jwt := ctx.GetHeader("Authorization")
	if jwt != "" {
		// Trim the "Bearer " prefix
		jwt = jwt[7:]
	}

	if jwt == "" {
		response.Data = false
		ctx.JSON(http.StatusOK, response)
		return
	}

	if _, err := c.authenticationService.VerifyToken(ctx, jwt); err != nil {
		response.Data = false
		ctx.JSON(http.StatusOK, response)
		return
	}

	claims := c.authenticationService.ExtractJWTData(ctx, jwt)

	if claims == nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	user := &models.UserModel{}

	if err := c.userSerivce.GetUserByID(ctx, claims.UserID, user); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if !c.authenticationService.CheckPassword(ctx, verifyPasswordStruct.Password, user.Password) {
		response.Data = false
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Data = true
	response.Message = config.GetMessageCode("VERIFY_SUCCESS")
	ctx.JSON(http.StatusOK, response)
}
