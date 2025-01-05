package services

import (
	"api/repositories"

	"github.com/gin-gonic/gin"
)

type BuildingService struct {
	buildingRepository *repositories.BuildingRepository
}

func NewBuildingService() *BuildingService {
	return &BuildingService{
		buildingRepository: repositories.NewBuildingRepository(),
	}
}

func (s *BuildingService) GetBuilding(ctx *gin.Context) {

}
