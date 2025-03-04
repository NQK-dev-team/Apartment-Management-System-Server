package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
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
