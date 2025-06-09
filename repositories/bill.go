package repositories

import (
	"api/config"
	"api/models"
	"errors"
	"time"

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

func (r *BillRepository) GetByContractId(ctx *gin.Context, bills *[]models.BillModel, contractID int64) error {
	if err := config.DB.Where("contract_id = ?", contractID).Preload("ExtraPayments").Find(bills).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	return nil
}

func (r *BillRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.BillModel{}).Where("id IN ?", id).UpdateColumns(models.BillModel{
		DefaultModel: models.DefaultModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}

	if err := tx.Set("isQuiet", true).Model(&models.ExtraPaymentModel{}).Where("bill_id IN ?", id).UpdateColumns(models.ExtraPaymentModel{
		DefaultModel: models.DefaultModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}

	return nil
}
