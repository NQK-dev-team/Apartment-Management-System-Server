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
