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
		_, err := m.authenticationService.VerifyToken(ctx, jwt)
		if err != nil {
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

	if response.JWTToken != jwt {
		ctx.Set("jwt", response.JWTToken)
	}

	ctx.Next()
}
