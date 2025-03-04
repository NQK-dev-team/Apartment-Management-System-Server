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

		newBuildingImageIDStart, err := s.buildingRepository.GetNewImageID(ctx)
		if err != nil {
			return err
		}

		newBuildingServiceIDStart, err := s.buildingRepository.GetNewServiceID(ctx)
		if err != nil {
			return err
		}

		newRoomID, err := s.roomRepository.GetNewID(ctx)
		if err != nil {
			return err
		}

		newRoomImageIDStart, err := s.roomRepository.GetNewImageID(ctx)
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

		for index, image := range building.Images {
			fileName, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/")
			if err != nil {
				return err
			}
			newBuilding.Images = append(newBuilding.Images, models.BuildingImageModel{
				BuildingID: newBuildingID,
				DefaultFileModel: models.DefaultFileModel{
					Path: fileName,
					No:   index + 1,
					ID:   newBuildingImageIDStart + int64(index),
				},
			})
			deleteImageList = append(deleteImageList, fileName)
		}

		for roomLoopIndex, val := range building.Rooms {
			newRoom := models.RoomModel{
				No:          val.No,
				Floor:       val.Floor,
				Status:      val.Status,
				Area:        val.Area,
				Description: val.Description,
				BuildingID:  newBuildingID,
				DefaultModel: models.DefaultModel{
					ID: newRoomID + int64(roomLoopIndex),
				},
			}

			roomNoStr := strconv.Itoa(int(val.No))

			for roomImageLoopIndex, image := range val.Images {
				fileName, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/"+"rooms/"+roomNoStr+"/")
				if err != nil {
					return err
				}
				newRoom.Images = append(newRoom.Images, models.RoomImageModel{
					RoomID: newRoomID,
					DefaultFileModel: models.DefaultFileModel{
						Path: fileName,
						No:   roomImageLoopIndex + 1,
						ID:   newRoomImageIDStart + int64(roomLoopIndex*(len(building.Rooms)-1)) + int64(roomImageLoopIndex),
					},
				})
			}
			newBuilding.Rooms = append(newBuilding.Rooms, newRoom)
		}

		for index, val := range building.Services {
			newBuilding.Services = append(newBuilding.Services, models.BuildingServiceModel{
				Name:       val.Name,
				Price:      val.Price,
				BuildingID: newBuildingID,
				DefaultModel: models.DefaultModel{
					ID: newBuildingServiceIDStart + int64(index),
				},
			})
		}

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

func (s *BuildingService) DeleteBuilding(ctx *gin.Context, id int64) error {
	deletedBuilding := &models.BuildingModel{
		DefaultModel: models.DefaultModel{
			ID: id,
		},
	}

	if err := s.buildingRepository.GetById(ctx, deletedBuilding, id); err != nil {
		return err
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}

	deletedBuilding.DeletedBy.Valid = true
	deletedBuilding.DeletedBy.Value = userID.(int64)

	for _, image := range deletedBuilding.Images {
		image.DeletedBy.Valid = true
		image.DeletedBy.Value = userID.(int64)
	}

	for _, service := range deletedBuilding.Services {
		service.DeletedBy.Valid = true
		service.DeletedBy.Value = userID.(int64)
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.buildingRepository.Update(ctx, deletedBuilding); err != nil {
			return err
		}

		if err := s.buildingRepository.Delete(ctx, deletedBuilding); err != nil {
			return err
		}

		return nil
	})
}
