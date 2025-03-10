package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
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
