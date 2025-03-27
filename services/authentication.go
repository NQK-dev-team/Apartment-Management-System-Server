package services

import (
	"api/config"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"database/sql"
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
		if err := s.refreshTokenRepository.Create(ctx, tx, &refreshToken); err != nil {
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

	token, _ := utils.ValidateJWTToken(jwtToken)
	claims := structs.JTWClaim{}

	if token == nil {
		return false, nil
	}

	utils.ExtractJWTClaim(token, &claims)

	user := models.UserModel{}

	if err := s.userRepository.GetByID(ctx, &user, claims.UserID); err != nil {
		return false, err
	}

	return user.IsOwner == claims.IsOwner && user.IsManager == claims.IsManager && user.IsCustomer == claims.IsCustomer, nil
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

func (s *AuthenticationService) CheckResetPasswordToken(ctx *gin.Context, token string, email string) (bool, error) {
	passwordResetToken := []models.PasswordResetTokenModel{}
	if err := s.passwordResetTokenRepository.GetByEmail(ctx, email, &passwordResetToken); err != nil {
		return false, err
	}

	if len(passwordResetToken) == 0 {
		return false, nil
	}

	if !utils.CompareHashString(passwordResetToken[0].Token, token) {
		return false, nil
	}

	if passwordResetToken[0].ExpiresAt.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func (s *AuthenticationService) CheckEmailVerifyToken(ctx *gin.Context, verifyEmailToken structs.VerifyEmailToken) (bool, error) {
	emailVerifyToken := []models.EmailVerifyTokenModel{}
	if err := s.emailVerifyTokenRepository.GetByEmail(ctx, verifyEmailToken.Email, &emailVerifyToken); err != nil {
		return false, err
	}

	if len(emailVerifyToken) == 0 {
		return false, nil
	}

	if !utils.CompareHashString(emailVerifyToken[0].Token, verifyEmailToken.Token) {
		return false, nil
	}

	if emailVerifyToken[0].ExpiresAt.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func (s *AuthenticationService) VerifyEmail(ctx *gin.Context, verifyEmailToken structs.VerifyEmailToken) error {
	var user = &models.UserModel{}
	s.userRepository.GetByEmail(ctx, user, verifyEmailToken.Email)

	user.EmailVerifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.Update(ctx, tx, user); err != nil {
			return err
		}
		if err := s.emailVerifyTokenRepository.Delete(ctx, tx, verifyEmailToken.Email); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *AuthenticationService) DeleteRefreshToken(ctx *gin.Context, userID int64) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.refreshTokenRepository.Delete(ctx, tx, userID); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *AuthenticationService) DeletePasswordResetToken(ctx *gin.Context, email string) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.passwordResetTokenRepository.Delete(ctx, tx, email); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *AuthenticationService) CheckPassword(ctx *gin.Context, providedPassword string, userPassword string) bool {
	return utils.CompareHashPassword(userPassword, providedPassword)
}
