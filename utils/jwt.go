package utils

import (
	"api/config"
	"api/structs"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(claim structs.JWTPayload) (string, error) {
	secretKey := config.GetEnv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET_KEY environment variable is not set")
	}
	expireTimeConfig := config.GetEnv("JWT_EXPIRE_TIME")
	if expireTimeConfig == "" {
		return "", errors.New("JWT_EXPIRE_TIME environment variable is not set")
	}
	expireTime, err := strconv.Atoi(expireTimeConfig)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID":       claim.UserID,
		"fullName":     claim.FullName,
		"imagePath":    claim.ImagePath,
		"isCustomer":   claim.IsCustomer,
		"isManager":    claim.IsManager,
		"isOwner":      claim.IsOwner,
		"userNo":       claim.UserNo,
		"ticketByPass": claim.TicketByPass,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(time.Second * time.Duration(expireTime)).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func ValidateJWTToken(tokenString string) (*jwt.Token, error) {
	secretKey := config.GetEnv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY environment variable is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractJWTClaim(token *jwt.Token, outputClaim *structs.JTWClaim) {
	claims := token.Claims.(jwt.MapClaims)
	outputClaim.UserID = int64(claims["userID"].(float64))
	outputClaim.FullName = claims["fullName"].(string)
	outputClaim.ImagePath = claims["imagePath"].(string)
	outputClaim.IsCustomer = claims["isCustomer"].(bool)
	outputClaim.IsManager = claims["isManager"].(bool)
	outputClaim.IsOwner = claims["isOwner"].(bool)
	outputClaim.UserNo = claims["userNo"].(string)
	outputClaim.TicketByPass = claims["ticketByPass"].(bool)
	outputClaim.ServiceToken = token.Raw
	outputClaim.IAT = int64(claims["iat"].(float64))
	outputClaim.EXP = int64(claims["exp"].(float64))
}

func GetRoleString(claims *structs.JTWClaim) string {
	var str = ""
	if claims.IsOwner {
		str += "1"
	} else {
		str += "0"
	}

	if claims.IsManager {
		str += "1"
	} else {
		str += "0"
	}

	if claims.IsCustomer {
		str += "1"
	} else {
		str += "0"
	}

	return str
}
