package middlewares

import (
	"api/config"

	"github.com/gin-gonic/gin"
)

type AuthorizationMiddleware struct {
}

func NewAuthorizationMiddleware() *AuthorizationMiddleware {
	return &AuthorizationMiddleware{}
}

func (m *AuthorizationMiddleware) AuthOwnerMiddleware(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	response := config.NewDataResponse(ctx)

	if role != "110" {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.AbortWithStatusJSON(403, response)
		return
	}

	ctx.Next()
}

func (m *AuthorizationMiddleware) AuthManagerMiddleware(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	response := config.NewDataResponse(ctx)

	if role != "110" && role != "010" {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.AbortWithStatusJSON(403, response)
		return
	}

	ctx.Next()
}

func (m *AuthorizationMiddleware) AuthCustomerMiddleware(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	response := config.NewDataResponse(ctx)

	if role != "001" {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.AbortWithStatusJSON(403, response)
		return
	}

	ctx.Next()
}
