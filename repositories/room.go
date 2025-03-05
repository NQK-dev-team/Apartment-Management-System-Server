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
	if err := config.DB.Preload("Building").Where("building_id = ?", buildingID).Find(room).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) GetNewID(ctx *gin.Context) (int64, error) {
	lastestRoom := models.RoomModel{}
	if err := config.DB.Order("id desc").First(&lastestRoom).Error; err != nil {
		return 0, err
	}
	return lastestRoom.ID + 1, nil
}
