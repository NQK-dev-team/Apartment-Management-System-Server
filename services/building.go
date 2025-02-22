package services

import (
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"

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

func (s *BuildingService) GetBuilding(ctx *gin.Context, building *[]models.BuildingModel) (bool, error) {
	role, exists := ctx.Get("role")

	if !exists {
		return false, nil
	}

	if role.(string) == constants.Roles.Manager {
		jwt, exists := ctx.Get("jwt")

		if !exists {
			return false, nil
		}

		token, err := utils.ValidateJWTToken(jwt.(string))

		if err != nil {
			return true, err
		}

		claim := &structs.JTWClaim{}

		utils.ExtractJWTClaim(token, claim)

		return true, s.buildingRepository.GetBuildingBaseOnSchedule(ctx, building, claim.UserID)
	}

	return true, s.buildingRepository.Get(ctx, building)
}

func (s *BuildingService) GetBuildingRoom(ctx *gin.Context, buildingID int64, room *[]models.RoomModel) error {
	return s.roomRepository.GetBuildingRoom(ctx, buildingID, room)
}
