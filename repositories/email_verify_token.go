package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailVerifyTokenRepository struct{}

func NewEmailVerifyTokenRepository() *EmailVerifyTokenRepository {
	return &EmailVerifyTokenRepository{}
}

func (r *EmailVerifyTokenRepository) Create(ctx *gin.Context, verifyToken *models.EmailVerifyTokenModel) error {
	if err := config.DB.Create(verifyToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *EmailVerifyTokenRepository) GetByEmail(ctx *gin.Context, email string, tokens *[]models.EmailVerifyTokenModel) error {
	if err := config.DB.Where("email = ?", email).Order("created_at DESC").Find(tokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *EmailVerifyTokenRepository) Delete(ctx *gin.Context, tx *gorm.DB, email string) error {
	if err := tx.Where("email = ?", email).Delete(&models.EmailVerifyTokenModel{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}
