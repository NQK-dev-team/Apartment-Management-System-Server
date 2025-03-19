package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"errors"
	"path/filepath"
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

func (s *BuildingService) CreateBuilding(ctx *gin.Context, building *structs.NewBuilding) error {
	deleteImageList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		newBuilding := &models.BuildingModel{}
		newBuilding.Name = building.Name
		newBuilding.Address = building.Address
		newBuilding.TotalFloor = building.TotalFloor
		newBuilding.TotalRoom = building.TotalRoom

		if err := s.buildingRepository.Create(ctx, tx, newBuilding); err != nil {
			return err
		}

		newBuildingIDStr := strconv.Itoa(int(newBuilding.ID))
		newImages := []models.BuildingImageModel{}
		for index, image := range building.Images {
			filePath, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/")
			if err != nil {
				return err
			}
			newImages = append(newImages, models.BuildingImageModel{
				BuildingID: newBuilding.ID,
				DefaultFileModel: models.DefaultFileModel{
					Path:  filePath,
					No:    index + 1,
					Title: filepath.Base(filePath),
				},
			})
			deleteImageList = append(deleteImageList, filePath)
		}

		if err := s.buildingRepository.AddImage(ctx, tx, &newImages); err != nil {
			return err
		}

		services := []models.BuildingServiceModel{}
		for _, val := range building.Services {
			services = append(services, models.BuildingServiceModel{
				Name:       val.Name,
				Price:      val.Price,
				BuildingID: newBuilding.ID,
			})
		}

		if len(services) > 0 {
			if err := s.buildingRepository.AddServices(ctx, tx, &services); err != nil {
				return err
			}
		}
		schedules := []models.ManagerScheduleModel{}
		for _, val := range building.Schedules {
			startDate := utils.ParseTime(val.StartDate)
			endDate := utils.StringToNullTime(val.EndDate)
			schedules = append(schedules, models.ManagerScheduleModel{
				BuildingID: newBuilding.ID,
				ManagerID:  val.ManagerID,
				StartDate:  startDate,
				EndDate:    endDate,
			})
		}
		if len(schedules) > 0 {
			if err := s.managerScheduleRepository.Create(ctx, tx, &schedules); err != nil {
				return err
			}
		}

		// for roomLoopIndex, val := range building.Rooms {
		// 	newRoom := models.RoomModel{
		// 		No:          val.No,
		// 		Floor:       val.Floor,
		// 		Status:      val.Status,
		// 		Area:        val.Area,
		// 		Description: val.Description,
		// 		BuildingID:  newBuildingID,
		// 		DefaultModel: models.DefaultModel{
		// 			ID: newRoomID + int64(roomLoopIndex),
		// 		},
		// 	}

		// 	roomNoStr := strconv.Itoa(int(val.No))

		// 	for roomImageLoopIndex, image := range val.Images {
		// 		filePath, err := utils.StoreFile(image, "images/buildings/"+newBuildingIDStr+"/"+"rooms/"+roomNoStr+"/")
		// 		if err != nil {
		// 			return err
		// 		}
		// 		newRoom.Images = append(newRoom.Images, models.RoomImageModel{
		// 			RoomID: newRoomID,
		// 			DefaultFileModel: models.DefaultFileModel{
		// 				Path:  filePath,
		// 				No:    roomImageLoopIndex + 1,
		// 				ID:    newRoomImageIDStart + int64(roomLoopIndex*(len(building.Rooms)-1)) + int64(roomImageLoopIndex),
		// 				Title: filepath.Base(filePath),
		// 			},
		// 		})
		// 		deleteImageList = append(deleteImageList, filePath)
		// 	}
		// 	newBuilding.Rooms = append(newBuilding.Rooms, newRoom)
		// }

		return errors.New("fake error")
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

func (s *BuildingService) GetBuildingSchedule(ctx *gin.Context, buildingID int64, schedules *[]models.ManagerScheduleModel) (bool, error) {
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

		return true, s.buildingRepository.GetManagerBuildingSchedule(ctx, buildingID, schedules, claim.UserID)
	}

	return true, s.buildingRepository.GetBuildingSchedule(ctx, buildingID, schedules)
}

func (s *BuildingService) UpdateBuilding(ctx *gin.Context, building *structs.EditBuilding) error {
	// role, exists := ctx.Get("role")

	// if !exists {
	// 	return errors.New("role not found")
	// }

	// buildingIDStr := strconv.Itoa(int(building.ID))

	deleteImageList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// newBuildingData := &models.BuildingModel{}

		// if err := s.buildingRepository.GetById(ctx, newBuildingData, building.ID); err != nil {
		// 	return err
		// }

		// if role.(string) == constants.Roles.Owner {
		// 	newBuildingData.Name = building.Name
		// 	newBuildingData.Address = building.Address
		// 	newBuildingData.TotalFloor = building.TotalFloor
		// 	newBuildingData.TotalRoom = len(building.Rooms) + len(building.NewRooms)

		// 	// newImageNo, err := s.buildingRepository.GetNewImageNo(ctx, building.ID)

		// 	// if err != nil {
		// 	// 	return err
		// 	// }

		// 	if err := s.buildingRepository.Update(ctx, newBuildingData); err != nil {
		// 		return err
		// 	}

		// 	if len(building.DeletedBuildingImages) > 0 {
		// 		if err := s.buildingRepository.DeleteImages(ctx, building.DeletedBuildingImages); err != nil {
		// 			return err
		// 		}
		// 	}

		// 	if len(building.DeletedSchedules) > 0 {
		// 		if err := s.managerScheduleRepository.Delete(ctx, building.DeletedSchedules); err != nil {
		// 			return err
		// 		}
		// 	}

		// 	if len(building.DeletedRooms) > 0 {
		// 		if err := s.roomRepository.Delete(ctx, building.DeletedRooms); err != nil {
		// 			return err
		// 		}
		// 	}

		// 	newImage := []models.BuildingImageModel{}
		// 	for _, image := range building.NewBuildingImages {
		// 		filePath, err := utils.StoreFile(image, "images/buildings/"+buildingIDStr+"/")
		// 		if err != nil {
		// 			return err
		// 		}
		// 		deleteImageList = append(deleteImageList, filePath)
		// 		newImage = append(newImage, models.BuildingImageModel{
		// 			BuildingID: building.ID,
		// 			DefaultFileModel: models.DefaultFileModel{
		// 				Title: filepath.Base(filePath),
		// 				Path:  filePath,
		// 			},
		// 			// Title: filepath.Base(filePath),
		// 			// Path:  filePath,
		// 			// No:         newImageNo + index,
		// 		})
		// 	}
		// 	if err := s.buildingRepository.AddImage(ctx, tx, &newImage); err != nil {
		// 		return err
		// 	}

		// 	schedules := []structs.NewBuildingSchedule{}
		// 	for _, schedule := range building.NewSchedules {
		// 		schedules = append(schedules, structs.NewBuildingSchedule{
		// 			BuildingID: building.ID,
		// 			ManagerID:  schedule.ManagerID,
		// 			StartDate:  utils.ParseTime(schedule.StartDate),
		// 			EndDate:    utils.StringToNullTime(schedule.EndDate),
		// 		})
		// 	}
		// 	if err := s.buildingRepository.AddBuilingSchedule(ctx, &schedules); err != nil {
		// 		return err
		// 	}

		// 	rooms := []structs.NewBuildingRoom{}
		// 	for _, room := range building.NewRooms {
		// 		images := []structs.NewRoomImage{}
		// 		roomNoStr := strconv.Itoa(int(room.No))
		// 		for _, image := range room.Images {
		// 			filePath, err := utils.StoreFile(image, "images/buildings/"+buildingIDStr+"/"+"rooms/"+roomNoStr+"/")
		// 			if err != nil {
		// 				return err
		// 			}
		// 			images = append(images, structs.NewRoomImage{
		// 				// RoomID: newRoomID,
		// 				Path:  filePath,
		// 				Title: filepath.Base(filePath),
		// 			})
		// 			deleteImageList = append(deleteImageList, filePath)
		// 		}

		// 		rooms = append(rooms, structs.NewBuildingRoom{
		// 			No:          room.No,
		// 			Images:      images,
		// 			Floor:       room.Floor,
		// 			Status:      room.Status,
		// 			Area:        room.Area,
		// 			Description: room.Description,
		// 			BuildingID:  building.ID,
		// 		})
		// 	}
		// 	if err := s.roomRepository.CreateRoom(ctx, &rooms); err != nil {
		// 		return err
		// 	}
		// }

		// if len(building.DeletedServices) > 0 {
		// 	if err := s.buildingRepository.DeleteServices(ctx, building.DeletedServices); err != nil {
		// 		return err
		// 	}
		// }

		// if len(building.NewServices) > 0 {
		// 	services := []structs.NewBuildingService{}
		// 	for _, service := range building.NewServices {
		// 		services = append(services, structs.NewBuildingService{
		// 			Name:       service.Name,
		// 			Price:      service.Price,
		// 			BuildingID: building.ID,
		// 		})
		// 	}
		// 	if err := s.buildingRepository.AddBuildingService(ctx, &services); err != nil {
		// 		return err
		// 	}
		// }

		return errors.New("fake error")
		// return nil
	})

	if err != nil {
		for _, path := range deleteImageList {
			utils.RemoveFile(path)
		}
		return err
	}

	return nil
}
