package repositories

import (
	"api/config"
	"api/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomRepository struct{}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (r *RoomRepository) GetBuildingRoom(ctx *gin.Context, buildingID int64, room *[]models.RoomModel) error {
	if err := config.DB.Where("building_id = ?", buildingID).Find(room).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) GetById(ctx *gin.Context, room *models.RoomModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("Images").Preload("Contracts").First(room).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) GetNewID(ctx *gin.Context) (int64, error) {
	lastestRoom := models.RoomModel{}
	if err := config.DB.Order("id desc").Unscoped().First(&lastestRoom).Error; err != nil {
		return 0, err
	}
	return lastestRoom.ID + 1, nil
}

func (r *RoomRepository) GetNewImageID(ctx *gin.Context) (int64, error) {
	lastestImage := models.RoomImageModel{}
	if err := config.DB.Order("id desc").Unscoped().First(&lastestImage).Error; err != nil {
		return 0, err
	}
	return lastestImage.ID + 1, nil
}

func (r *RoomRepository) Delete(ctx *gin.Context, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := config.DB.Set("isQuiet", true).Model(&models.RoomModel{}).Where("id in ?", id).UpdateColumns(models.RoomModel{
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

	if err := config.DB.Set("isQuiet", true).Model(&models.RoomImageModel{}).Where("room_id in ?", id).UpdateColumns(models.RoomImageModel{
		DefaultFileModel: models.DefaultFileModel{
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
