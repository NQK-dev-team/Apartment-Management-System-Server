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

func (r *BillRepository) Get(ctx *gin.Context, bill *[]models.BillModel) error {
	if err := config.DB.Find(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetBillBaseOnSchedule(ctx *gin.Context, bill *[]models.BillModel, userID int64) error {
	if err := config.DB.Model(&models.BillModel{}).Select("bill.*").Joins("JOIN manager_schedule ON manager_schedule.bill_id = bill.id").Where("manager_schedule.start_date <= now() AND manager_schedule.end_date >= now() AND manager_schedule.manager_id = ?", userID).Scan(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetById(ctx *gin.Context, bill *models.BillModel, id int64) error {
	if err := config.DB.Where("id = ?", id).First(bill).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *BillRepository) Create(ctx *gin.Context, bill *models.BillModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Create(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) Update(ctx *gin.Context, bill *models.BillModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Save(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) Delete(ctx *gin.Context, bill *models.BillModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Delete(bill).Error; err != nil {
		return err
	}
	return nil
}