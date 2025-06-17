package middlewares

import (
	"api/config"
	"api/services"
	"api/utils"

	"github.com/gin-gonic/gin"
)

type AuthenticationMiddleware struct {
	authenticationService *services.AuthenticationService
}

func NewAuthenticationMiddleware() *AuthenticationMiddleware {
	authService := services.NewAuthenticationService()
	return &AuthenticationMiddleware{authenticationService: authService}
}

func (m *AuthenticationMiddleware) AuthMiddleware(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	// Get the JWT token from the header
	jwt := ctx.GetHeader("Authorization")
	if jwt != "" {
		// Trim the "Bearer " prefix
		jwt = jwt[7:]
	}
	// Get the refresh token from the header
	refreshToken := ctx.GetHeader("X-Refresh-Token")

	if jwt == "" {
		if refreshToken == "" {
			response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
			ctx.AbortWithStatusJSON(401, response)
			return
		} else {
			newJwt, err := m.authenticationService.GetNewToken(ctx, refreshToken)

			if err != nil {
				response.Message = config.GetMessageCode("TOKEN_REFRESH_FAILED")
				ctx.AbortWithStatusJSON(401, response)
				return
			}

			response.JWTToken = newJwt
		}
	} else {
		isValid, err := m.authenticationService.VerifyToken(ctx, jwt)
		if err != nil || !isValid {
			if refreshToken != "" {
				newJwt, err := m.authenticationService.GetNewToken(ctx, refreshToken)

				if err != nil {
					response.Message = config.GetMessageCode("TOKEN_REFRESH_FAILED")
					ctx.AbortWithStatusJSON(401, response)
					return
				}

				response.JWTToken = newJwt
			} else {
				response.Message = config.GetMessageCode("TOKEN_VERIFY_FAILED")
				ctx.AbortWithStatusJSON(401, response)
				return
			}
		} else {
			response.JWTToken = jwt
		}
	}

	claims := m.authenticationService.ExtractJWTData(ctx, response.JWTToken)

	if claims == nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.AbortWithStatusJSON(500, response)
		return
	}

	ctx.Set("role", utils.GetRoleString(claims))

	// Get the refresh token
	refreshTokenRecord := &models.RefreshTokenModel{}
	if err := m.authenticationService.GetRefreshToken(ctx, refreshTokenRecord, claims.UserID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// Check if the refresh token expiration date is less than 1 day from now
	if refreshTokenRecord.ExpiresAt.Before(time.Now().AddDate(0, 0, 1)) {
		// Delete the old refresh token
		if err := m.authenticationService.DeleteRefreshToken(ctx, claims.UserID); err != nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		// Generate a new refresh token
		newRefreshToken, err := m.authenticationService.CreateRefreshToken(ctx, claims.UserID)
		if err != nil {
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		ctx.Set("refreshToken", newRefreshToken)
	}

	// if response.JWTToken != jwt {
	ctx.Set("jwt", response.JWTToken)
	// }

	ctx.Set("userID", claims.UserID)

	ctx.Next()
}
