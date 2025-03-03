package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (s *BuildingService) CreateBuilding(ctx *gin.Context, building *structs.NewBuilding) error {
	deleteImageList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		newBuildingID, err := s.buildingRepository.GetNewID(ctx)
		if err != nil {
			return err
		}

		newBuildingIDStr := strconv.Itoa(int(newBuildingID))

		newBuilding := &models.BuildingModel{}
		newBuilding.Name = building.Name
		newBuilding.ID = newBuildingID
		newBuilding.Address = building.Address
		newBuilding.TotalFloor = building.TotalFloor
		newBuilding.TotalRoom = building.TotalRoom

		for _, image := range building.Images {
			fileName, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/")
			if err != nil {
				return err
			}
			newBuilding.Images = append(newBuilding.Images, models.BuildingImageModel{
				BuildingID: newBuildingID,
				DefaultFileModel: models.DefaultFileModel{
					Path: fileName,
				},
			})
			deleteImageList = append(deleteImageList, fileName)
		}

		for index, val := range building.Rooms {
			newRoomID, err := s.roomRepository.GetNewID(ctx)
			if err != nil {
				return err
			}

			newRoomID = newRoomID + int64(index)

			newRoom := models.RoomModel{
				No:          val.No,
				Floor:       val.Floor,
				Status:      val.Status,
				Area:        val.Area,
				Description: val.Description,
				BuildingID:  newBuildingID,
				DefaultModel: models.DefaultModel{
					ID: newRoomID,
				},
			}

			roomNoStr := strconv.Itoa(int(val.No))

			for _, image := range val.Images {
				fileName, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/"+"rooms/"+roomNoStr+"/")
				if err != nil {
					return err
				}
				newRoom.Images = append(newRoom.Images, models.RoomImageModel{
					BuildingID: newBuildingID,
					RoomID:     newRoomID,
					DefaultFileModel: models.DefaultFileModel{
						Path: fileName,
					},
				})
			}
			newBuilding.Rooms = append(newBuilding.Rooms, newRoom)
		}

		for _, val := range building.Services {
			newBuilding.Services = append(newBuilding.Services, models.BuildingServiceModel{
				Name:       val.Name,
				Price:      val.Price,
				BuildingID: newBuildingID,
			})
		}

		ctx.Set("userID", utils.GetUserID(ctx))

		if err := s.buildingRepository.Create(ctx, newBuilding); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		for _, path := range deleteImageList {
			utils.RemoveFile(path)
		}
		return err
	}
	return nil
}

func (s *BuildingService) GetBuildingDetail(ctx *gin.Context, building *models.BuildingModel, id int64) error {
	return s.buildingRepository.GetById(ctx, building, id)
}
