package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadRepository struct{}

func NewUploadRepository() *UploadRepository {
	return &UploadRepository{}
}

func (r *UploadRepository) Get(ctx *gin.Context, uploads *[]models.UploadFileModel, uploadType int, isProcessed bool) error {
	query := config.DB.Model(&models.UploadFileModel{}).Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Select("id", "no", "first_name", "middle_name", "last_name", "is_owner", "is_manager", "is_customer")
	})

	query = query.Where("upload_type = ?", uploadType)

	if isProcessed {
		query = query.Where("process_result IS NOT NULL")
	} else {
		query = query.Where("process_result IS NULL")
	}

	if err := query.Find(uploads).Error; err != nil {
		return err
	}
	return nil
}

func (r *UploadRepository) Create(ctx *gin.Context, tx *gorm.DB, upload *models.UploadFileModel) error {
	if err := tx.Model(&models.UploadFileModel{}).Omit("ID").Create(upload).Error; err != nil {
		return err
	}
	return nil
}

func (r *UploadRepository) Update(ctx *gin.Context, tx *gorm.DB, upload *models.UploadFileModel) error {
	if err := tx.Model(&models.UploadFileModel{}).Where("id = ?", upload.ID).Updates(upload).Error; err != nil {
		return err
	}
	return nil
}

func(r*UploadRepository) GetByID(ctx *gin.Context, upload *models.UploadFileModel, uploadID int64) error {
	if err := config.DB.Model(&models.UploadFileModel{}).Where("id = ?", uploadID).Find(upload).Error; err != nil {
		return err
	}
	return nil
}