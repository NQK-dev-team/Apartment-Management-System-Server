package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
)

type RefreshTokenRepository struct {
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

// func (r *RefreshTokenRepository) Create(ctx *gin.Context, tx *gorm.DB, refreshToken *models.RefreshTokenModel) error {
// 	if err := tx.Create(refreshToken).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshTokenModel) error {
	if err := config.DBNoLog.Create(refreshToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) GetByUserID(ctx *gin.Context, refreshToken *models.RefreshTokenModel, userID int64) error {
	if err := config.DBNoLog.Where("user_id = ?", userID).Order("created_at DESC").Find(refreshToken).Error; err != nil {
		return err
	}

	return nil
}

// func (r *RefreshTokenRepository) Delete(ctx *gin.Context, tx *gorm.DB, userID int64) error {
// 	if err := tx.Where("user_id = ?", userID).Delete(&models.RefreshTokenModel{}).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *RefreshTokenRepository) Delete(userID int64) error {
	if err := config.DBNoLog.Where("user_id = ?", userID).Delete(&models.RefreshTokenModel{}).Error; err != nil {
		return err
	}

	return nil
}
