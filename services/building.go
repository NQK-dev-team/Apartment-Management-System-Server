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
	roomService               *RoomService
	buildingRepository        *repositories.BuildingRepository
	roomRepository            *repositories.RoomRepository
	managerScheduleRepository *repositories.ManagerScheduleRepository
}

func NewBuildingService() *BuildingService {
	return &BuildingService{
		roomService:               NewRoomService(),
		buildingRepository:        repositories.NewBuildingRepository(),
		roomRepository:            repositories.NewRoomRepository(),
		managerScheduleRepository: repositories.NewManagerScheduleRepository(),
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

func (s *BuildingService) GetBuildingService(ctx *gin.Context, buildingID int64, service *[]models.BuildingServiceModel) error {
	return s.buildingRepository.GetBuildingService(ctx, buildingID, service)
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
			filePath, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/")
			if err != nil {
				return err
			}
			newBuilding.Images = append(newBuilding.Images, models.BuildingImageModel{
				BuildingID: newBuildingID,
				DefaultFileModel: models.DefaultFileModel{
					Path: filePath,
					No:   index + 1,
					ID:   newBuildingImageIDStart + int64(index),
				},
			})
			deleteImageList = append(deleteImageList, filePath)
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
				filePath, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/"+"rooms/"+roomNoStr+"/")
				if err != nil {
					return err
				}
				newRoom.Images = append(newRoom.Images, models.RoomImageModel{
					RoomID: newRoomID,
					DefaultFileModel: models.DefaultFileModel{
						Path: filePath,
						No:   roomImageLoopIndex + 1,
						ID:   newRoomImageIDStart + int64(roomLoopIndex*(len(building.Rooms)-1)) + int64(roomImageLoopIndex),
					},
				})
				deleteImageList = append(deleteImageList, filePath)
			}
			newBuilding.Rooms = append(newBuilding.Rooms, newRoom)
		}

		if err := s.buildingRepository.Create(ctx, newBuilding); err != nil {
			return err
		}

		newScheduleIDStart, err := s.managerScheduleRepository.GetNewScheduleID(ctx)

		if err != nil {
			return err
		}

		schedules := []models.ManagerScheduleModel{}

		for index, val := range building.Schedules {
			startDate := utils.ParseTime(val.StartDate)
			endDate := utils.StringToNullTime(val.EndDate)
			schedules = append(schedules, models.ManagerScheduleModel{
				BuildingID: newBuildingID,
				ManagerID:  val.ManagerID,
				StartDate:  startDate,
				EndDate:    endDate,
				DefaultModel: models.DefaultModel{
					ID: newScheduleIDStart + int64(index),
				},
			})
		}

		if err := s.managerScheduleRepository.Create(ctx, &schedules); err != nil {
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

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.buildingRepository.Delete(ctx, []int64{deletedBuilding.ID}); err != nil {
			return err
		}

		roomIDs := []int64{}
		for _, room := range deletedBuilding.Rooms {
			roomIDs = append(roomIDs, room.ID)
		}

		if err := s.roomService.DeleteWithoutTransaction(ctx, roomIDs); err != nil {
			return err
		}

		return nil
	})
}

func (s *BuildingService) DeleteRooms(ctx *gin.Context, buildingID int64, roomIDs []int64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.roomService.DeleteWithoutTransaction(ctx, roomIDs); err != nil {
			return err
		}

		building := &models.BuildingModel{}
		if err := s.buildingRepository.GetById(ctx, building, buildingID); err != nil {
			return err
		}

		building.TotalRoom = building.TotalRoom - len(roomIDs)

		if err := s.buildingRepository.Update(ctx, building); err != nil {
			return err
		}

		return nil
	})
}

func (s *BuildingService) DeleteServices(ctx *gin.Context, serviceIDs []int64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return s.buildingRepository.DeleteServices(ctx, serviceIDs)
	})
}

func (s *BuildingService) AddService(ctx *gin.Context, service *models.BuildingServiceModel) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		newID, err := s.buildingRepository.GetNewServiceID(ctx)
		if err != nil {
			return err
		}

		service.ID = newID

		return s.buildingRepository.AddService(ctx, service)
	})
}

func (s *BuildingService) EditService(ctx *gin.Context, newServiceData *models.BuildingServiceModel) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		service := &models.BuildingServiceModel{}
		if err := s.buildingRepository.GetServiceByID(ctx, service, newServiceData.ID); err != nil {
			return err
		}

		service.Name = newServiceData.Name
		service.Price = newServiceData.Price

		return s.buildingRepository.EditService(ctx, service)
	})
}

func (s *BuildingService) AddRoom(ctx *gin.Context, buildingID int64, room *structs.NewRoom) error {
	deleteImageList := []string{}
	newRoom := &models.RoomModel{
		No:          room.No,
		Floor:       room.Floor,
		Status:      room.Status,
		Area:        room.Area,
		Description: room.Description,
		BuildingID:  buildingID,
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		newRoomID, err := s.roomRepository.GetNewID(ctx)
		if err != nil {
			return err
		}

		newRoom.ID = newRoomID

		newRoomImageIDStart, err := s.roomRepository.GetNewImageID(ctx)
		if err != nil {
			return err
		}

		for index, image := range room.Images {
			filePath, err := utils.StoreFile(image, "images/buildings/"+strconv.Itoa(int(buildingID))+"/"+"rooms/"+strconv.Itoa(int(room.No))+"/")
			if err != nil {
				return err
			}
			newRoom.Images = append(newRoom.Images, models.RoomImageModel{
				RoomID: newRoomID,
				DefaultFileModel: models.DefaultFileModel{
					Path: filePath,
					No:   index + 1,
					ID:   newRoomImageIDStart + int64(index),
				},
			})
			deleteImageList = append(deleteImageList, filePath)
		}

		if err := s.roomRepository.Create(ctx, newRoom); err != nil {
			return err
		}

		building := &models.BuildingModel{}
		if err := s.buildingRepository.GetById(ctx, building, buildingID); err != nil {
			return err
		}

		building.TotalRoom = building.TotalRoom + 1

		if err := s.buildingRepository.Update(ctx, building); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		for _, path := range deleteImageList {
			utils.RemoveFile(path)
		}
	}

	return nil
}

func (s *BuildingService) GetBuildingSchedule(ctx *gin.Context, buildingID int64, schedules *[]models.ManagerScheduleModel) error {
	return s.buildingRepository.GetBuildingSchedule(ctx, buildingID, schedules)
}
