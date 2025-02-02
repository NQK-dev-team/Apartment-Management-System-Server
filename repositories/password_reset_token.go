package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PasswordResetTokenRepository struct {
}

func NewPasswordResetTokenRepository() *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{}
}

func (r *PasswordResetTokenRepository) Create(ctx *gin.Context, passwordResetToken *models.PasswordResetTokenModel) error {
	if err := config.DB.Create(passwordResetToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *PasswordResetTokenRepository) GetByEmail(ctx *gin.Context, email string, tokens *[]models.PasswordResetTokenModel) error {
	if err := config.DB.Where("email = ?", email).Order("created_at DESC").Find(tokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *PasswordResetTokenRepository) Delete(ctx *gin.Context, email string) error {
	if err := config.DB.Where("email = ?", email).Delete(&models.PasswordResetTokenModel{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}
