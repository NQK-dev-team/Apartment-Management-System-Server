package middlewares

import (
	"api/config"
	"api/constants"

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

	if role != constants.Roles.Owner {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.AbortWithStatusJSON(403, response)
		return
	}

	ctx.Next()
}

func (m *AuthorizationMiddleware) AuthManagerMiddleware(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	response := config.NewDataResponse(ctx)

	if role != constants.Roles.Owner && role != constants.Roles.Manager {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.AbortWithStatusJSON(403, response)
		return
	}

	ctx.Next()
}

func (m *AuthorizationMiddleware) AuthCustomerMiddleware(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	response := config.NewDataResponse(ctx)

	if role != constants.Roles.Customer {
		response.Message = config.GetMessageCode("PERMISSION_DENIED")
		ctx.AbortWithStatusJSON(403, response)
		return
	}

	ctx.Next()
}
