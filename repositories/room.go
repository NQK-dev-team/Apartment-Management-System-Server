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

func (r *RoomRepository) GetBuildingRooms(ctx *gin.Context, buildingID string, rooms *[]models.RoomModel) error {
	if err := config.DB.Where("building_id = ?", buildingID).Find(rooms).Error; err != nil {
		return err
	}
	return nil
}
