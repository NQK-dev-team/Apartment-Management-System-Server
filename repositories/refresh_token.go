package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

func (r *RefreshTokenRepository) Create(ctx *gin.Context, tx *gorm.DB, refreshToken *models.RefreshTokenModel) error {
	if err := tx.Create(refreshToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) GetByUserID(ctx *gin.Context, refreshToken *models.RefreshTokenModel, userID int64) error {
	if err := config.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(refreshToken).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil
		// }
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) Delete(ctx *gin.Context, tx *gorm.DB, userID int64) error {
	if err := tx.Where("user_id = ?", userID).Delete(&models.RefreshTokenModel{}).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil
		// }
		return err
	}

	return nil
}
