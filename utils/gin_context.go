package utils

import (
	"api/structs"

	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) int64 {
	jwt, exists := ctx.Get("jwt")

	if !exists {
		return 0
	}

	token, err := ValidateJWTToken(jwt.(string))

	if err != nil {
		return 0
	}

	claim := &structs.JTWClaim{}

	ExtractJWTClaim(token, claim)

	return claim.UserID
}
