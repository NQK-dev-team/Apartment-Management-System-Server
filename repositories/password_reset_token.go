package repositories

import (
	"api/config"
	"api/models"

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
		return err
	}

	return nil
}

func (r *PasswordResetTokenRepository) Delete(ctx *gin.Context, tx *gorm.DB, email string) error {
	if err := tx.Where("email = ?", email).Delete(&models.PasswordResetTokenModel{}).Error; err != nil {
		return err
	}

	return nil
}
