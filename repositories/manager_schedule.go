package repositories

import (
	"api/config"
	"api/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ManagerScheduleRepository struct {
}

func NewManagerScheduleRepository() *ManagerScheduleRepository {
	return &ManagerScheduleRepository{}
}

func (r *ManagerScheduleRepository) Create(ctx *gin.Context, tx *gorm.DB, schedules *[]models.ManagerScheduleModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.ManagerScheduleModel{}).Omit("ID").Create(schedules).Error; err != nil {
		return err
	}
	return nil
}

func (r *ManagerScheduleRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.ManagerScheduleModel{}).Where("id in ?", id).UpdateColumns(models.ManagerScheduleModel{
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

func (r *ManagerScheduleRepository) GetByIDs(ctx *gin.Context, schedule *[]models.ManagerScheduleModel, IDs []int64) error {
	if err := config.DB.Where("id in ?", IDs).Find(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *ManagerScheduleRepository) Update(ctx *gin.Context, tx *gorm.DB, schedules *[]models.ManagerScheduleModel) error {
	userID := ctx.GetInt64("userID")
	// if err := tx.Set("userID", userID).Updates(schedules).Error; err != nil {
	// 	return err
	// }

	for _, schedule := range *schedules {
		if err := tx.Set("userID", userID).Model(&models.ManagerScheduleModel{}).Where("id = ?", schedule.ID).Updates(schedule).Error; err != nil {
			return err
		}
	}

	return nil
}
