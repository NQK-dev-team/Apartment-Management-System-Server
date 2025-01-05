package services

import (
	"api/config"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	userRepository               *repositories.UserRepository
	refreshTokenRepository       *repositories.RefreshTokenRepository
	emailVerifyTokenRepository   *repositories.EmailVerifyTokenRepository
	passwordResetTokenRepository *repositories.PasswordResetTokenRepository
}

func NewAuthenticationService() *AuthenticationService {
	userRepository := repositories.NewUserRepository()
	refreshTokenRepository := repositories.NewRefreshTokenRepository()
	emailVerifyTokenRepository := repositories.NewEmailVerifyTokenRepository()
	passwordResetTokenRepository := repositories.NewPasswordResetTokenRepository()
	return &AuthenticationService{
		userRepository:               userRepository,
		refreshTokenRepository:       refreshTokenRepository,
		emailVerifyTokenRepository:   emailVerifyTokenRepository,
		passwordResetTokenRepository: passwordResetTokenRepository,
	}
}

func (s *AuthenticationService) Login(ctx *gin.Context, email string, password string, remember bool) (string, string, error, bool) {
	user := models.UserModel{}
	err := s.userRepository.GetByEmail(ctx, &user, email)

	if err != nil {
		return "", "", err, true
	}

	if user.ID == 0 || !utils.CompareHashPassword(user.Password, password) {
		return "", "", nil, true
	}

	if !user.EmailVerifiedAt.Valid {
		return "", "", nil, false
	}

	jwtPayload := structs.JWTPayload{}
	jwtPayload.UserID = user.ID
	jwtPayload.ImagePath = user.ProfileFilePath
	jwtPayload.IsCustomer = user.IsCustomer
	jwtPayload.IsManager = user.IsManager
	jwtPayload.IsOwner = user.IsOwner

	if user.MiddleName != "" {
		jwtPayload.FullName = user.FirstName + " " + user.MiddleName + " " + user.LastName
	} else {
		jwtPayload.FullName = user.FirstName + " " + user.LastName
	}

	// Create JWT token
	jwtToken, err := utils.GenerateJWTToken(jwtPayload)
	if err != nil {
		return "", "", err, true
	}

	var refreshToken = ""

	if remember {
		refreshToken, _ = s.CreateRefreshToken(ctx, user.ID)
	}

	return jwtToken, refreshToken, nil, true
}

func (s *AuthenticationService) CreateRefreshToken(ctx *gin.Context, userID int64) (string, error) {
	token, err := utils.GenerateString(64)
	if err != nil {
		return "", err
	}

	token = fmt.Sprintf("%d.%s", userID, token)

	hashedToken, err := utils.HashString(token)

	if err != nil {
		return "", err
	}

	refreshToken := models.RefreshTokenModel{
		Token:  hashedToken,
		UserID: userID,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.refreshTokenRepository.Create(ctx, &refreshToken); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthenticationService) GetNewToken(ctx *gin.Context, refreshToken string) (string, error) {
	arr := strings.Split(refreshToken, ".")
	userID, err := strconv.ParseInt(arr[0], 10, 64)

	if err != nil {
		return "", err
	}

	retrievedRefreshToken := models.RefreshTokenModel{}

	if err := s.refreshTokenRepository.GetByUserID(ctx, &retrievedRefreshToken, userID); err != nil {
		return "", err
	}

	if retrievedRefreshToken.ExpiresAt.Before(time.Now()) {
		return "", errors.New("refresh token has expired")
	}

	if !utils.CompareHashString(retrievedRefreshToken.Token, refreshToken) {
		return "", nil
	}

	user := models.UserModel{}

	if err := s.userRepository.GetByID(ctx, &user, userID); err != nil {
		return "", err
	}

	jwtPayload := structs.JWTPayload{}
	jwtPayload.UserID = user.ID
	jwtPayload.ImagePath = user.ProfileFilePath
	jwtPayload.IsCustomer = user.IsCustomer
	jwtPayload.IsManager = user.IsManager
	jwtPayload.IsOwner = user.IsOwner

	if user.MiddleName != "" {
		jwtPayload.FullName = user.FirstName + " " + user.MiddleName + " " + user.LastName
	} else {
		jwtPayload.FullName = user.FirstName + " " + user.LastName
	}

	// Create JWT token
	jwtToken, err := utils.GenerateJWTToken(jwtPayload)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s *AuthenticationService) VerifyToken(ctx *gin.Context, jwtToken string) (bool, error) {
	_, err := utils.ValidateJWTToken(jwtToken)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AuthenticationService) ExtractJWTData(ctx *gin.Context, jwt string) *structs.JTWClaim {
	token, _ := utils.ValidateJWTToken(jwt)
	claims := structs.JTWClaim{}

	if token == nil {
		return nil
	}

	utils.ExtractJWTClaim(token, &claims)

	return &claims
}

func (s *AuthenticationService) Logout(ctx *gin.Context) error {
	return nil
}
