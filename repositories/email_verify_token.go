package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
)

type EmailVerifyTokenRepository struct{}

func NewEmailVerifyTokenRepository() *EmailVerifyTokenRepository {
	return &EmailVerifyTokenRepository{}
}

func (r *EmailVerifyTokenRepository) Create(ctx *gin.Context, verifyToken *models.EmailVerifyTokenModel) error {
	if err := config.DBNoLog.Create(verifyToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *EmailVerifyTokenRepository) GetByEmail(ctx *gin.Context, email string, tokens *[]models.EmailVerifyTokenModel) error {
	if err := config.DBNoLog.Where("email = ?", email).Order("created_at DESC").Find(tokens).Error; err != nil {
		return err
	}

	return nil
}

// func (r *EmailVerifyTokenRepository) Delete(ctx *gin.Context, tx *gorm.DB, email string) error {
// 	if err := tx.Where("email = ?", email).Delete(&models.EmailVerifyTokenModel{}).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *EmailVerifyTokenRepository) Delete(email string) error {
	if err := config.DBNoLog.Where("email = ?", email).Delete(&models.EmailVerifyTokenModel{}).Error; err != nil {
		return err
	}

	return nil
}
