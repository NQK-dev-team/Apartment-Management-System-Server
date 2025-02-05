package services

import (
	"api/models"
	"api/repositories"

	"github.com/gin-gonic/gin"
)

type BuildingService struct {
	buildingRepository *repositories.BuildingRepository
	roomRepository     *repositories.RoomRepository
}

func NewBuildingService() *BuildingService {
	return &BuildingService{
		buildingRepository: repositories.NewBuildingRepository(),
		roomRepository:     repositories.NewRoomRepository(),
	}
}

func (s *BuildingService) GetBuilding(ctx *gin.Context, buildings *[]models.BuildingModel) error {
	return s.buildingRepository.Get(ctx, buildings)
}

func (s *BuildingService) GetBuildingRooms(ctx *gin.Context, buildingID string, rooms *[]models.RoomModel) error {
	return s.roomRepository.GetBuildingRooms(ctx, buildingID, rooms)
}
