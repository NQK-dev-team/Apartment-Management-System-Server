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

func (r *ManagerScheduleRepository) GetNewScheduleID(ctx *gin.Context) (int64, error) {
	lastestSchedule := models.ManagerScheduleModel{}
	if err := config.DB.Order("id desc").Unscoped().First(&lastestSchedule).Error; err != nil {
		return 0, err
	}
	return lastestSchedule.ID + 1, nil
}

func (r *ManagerScheduleRepository) Create(ctx *gin.Context, schedule *[]models.ManagerScheduleModel) error {
	userID := ctx.GetInt64("userID")
	if err := config.DB.Set("userID", userID).Create(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *ManagerScheduleRepository) Delete(ctx *gin.Context, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := config.DB.Set("isQuiet", true).Model(&models.ManagerScheduleModel{}).Where("id in ?", id).UpdateColumns(models.ManagerScheduleModel{
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
