package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BillRepository struct {
}

func NewBillRepository() *BillRepository {
	return &BillRepository{}
}

func (r *BillRepository) GetById(ctx *gin.Context, bill *models.BillModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("ExtraPayments").First(bill).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *BillRepository) QuietUpdate(ctx *gin.Context, bill *models.BillModel) error {
	if err := config.DB.Set("isQuiet", true).Session(&gorm.Session{FullSaveAssociations: true}).Save(bill).Error; err != nil {
		return err
	}
	return nil
}
