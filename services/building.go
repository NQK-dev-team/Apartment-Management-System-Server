package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
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

func NewBuildingService(noInitRoomService bool) *BuildingService {
	if noInitRoomService {
		return &BuildingService{
			roomService:               nil, // No RoomService initialized
			buildingRepository:        repositories.NewBuildingRepository(),
			roomRepository:            repositories.NewRoomRepository(),
			managerScheduleRepository: repositories.NewManagerScheduleRepository(),
		}
	}

	return &BuildingService{
		roomService:               NewRoomService(),
		buildingRepository:        repositories.NewBuildingRepository(),
		roomRepository:            repositories.NewRoomRepository(),
		managerScheduleRepository: repositories.NewManagerScheduleRepository(),
	}
}

func (s *BuildingService) GetBuilding(ctx *gin.Context, building *[]models.BuildingModel, getAll bool) (bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager && !getAll {
		return true, s.buildingRepository.GetBuildingBaseOnSchedule(ctx, building, ctx.GetInt64("userID"))
	}

	return true, s.buildingRepository.Get(ctx, building)
}

func (s *BuildingService) CreateBuilding(ctx *gin.Context, building *structs.NewBuilding, newBuildingID *int64) error {
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

		*newBuildingID = newBuilding.ID

		newBuildingIDStr := strconv.Itoa(int(newBuilding.ID))
		newImages := []models.BuildingImageModel{}
		for index, image := range building.Images {
			filePath, err := utils.StoreFile(image, constants.GetBuildingImageURL("images", newBuildingIDStr, ""))
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
			startDate, err := utils.ParseTime(val.StartDate)
			if err != nil {
				return err
			}
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

		rooms := []models.RoomModel{}
		for _, room := range building.Rooms {
			rooms = append(rooms, models.RoomModel{
				No:          room.No,
				Floor:       room.Floor,
				Status:      room.Status,
				Area:        room.Area,
				Description: room.Description,
				BuildingID:  newBuilding.ID,
			})
		}

		if len(rooms) > 0 {
			if err := s.roomRepository.Create(ctx, tx, &rooms); err != nil {
				return err
			}

			roomImages := []models.RoomImageModel{}
			for _, room := range building.Rooms {
				var targetRoom models.RoomModel
				for _, val := range rooms {
					if val.No == room.No {
						targetRoom = val
					}
				}
				roomNoStr := strconv.Itoa(int(targetRoom.No))

				for index, image := range room.Images {
					filePath, err := utils.StoreFile(image, constants.GetRoomImageURL("images", newBuildingIDStr, roomNoStr, ""))
					if err != nil {
						return err
					}
					roomImages = append(roomImages, models.RoomImageModel{
						RoomID: targetRoom.ID,
						DefaultFileModel: models.DefaultFileModel{
							Path:  filePath,
							No:    index + 1,
							Title: filepath.Base(filePath),
						},
					})
					deleteImageList = append(deleteImageList, filePath)
				}
			}

			if err := s.roomRepository.CreateImage(ctx, tx, &roomImages); err != nil {
				return err
			}
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
		if err := s.buildingRepository.Delete(ctx, tx, []int64{deletedBuilding.ID}); err != nil {
			return err
		}

		// roomIDs := []int64{}
		// for _, room := range deletedBuilding.Rooms {
		// 	roomIDs = append(roomIDs, room.ID)
		// }

		// if err := s.roomService.DeleteWithoutTransaction(ctx, tx, roomIDs); err != nil {
		// 	return err
		// }

		return nil
	})
}

func (s *BuildingService) GetBuildingSchedule(ctx *gin.Context, buildingID int64, schedules *[]models.ManagerScheduleModel) (bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager {
		return true, s.buildingRepository.GetManagerBuildingSchedule(ctx, buildingID, schedules, ctx.GetInt64("userID"))
	}

	return true, s.buildingRepository.GetBuildingSchedule(ctx, buildingID, schedules)
}

func (s *BuildingService) GetBuildingRoom(ctx *gin.Context, buildingID int64, rooms *[]models.RoomModel) (bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager {
		building := &[]models.BuildingModel{}

		if err := s.buildingRepository.GetBuildingBaseOnSchedule(ctx, building, ctx.GetInt64("userID")); err != nil {
			return false, err
		}

		if len(*building) == 0 {
			return false, nil
		}

		for _, val := range *building {
			if val.ID == buildingID {
				return true, s.buildingRepository.GetBuildingRoom(ctx, buildingID, rooms)
			}
		}
		return false, nil
	}

	return true, s.buildingRepository.GetBuildingRoom(ctx, buildingID, rooms)
}

func (s *BuildingService) UpdateBuilding(ctx *gin.Context, building *structs.EditBuilding) error {
	role := ctx.GetString("role")

	buildingIDStr := strconv.Itoa(int(building.ID))

	deleteImageList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		newBuildingData := &models.BuildingModel{}

		if err := s.buildingRepository.GetById(ctx, newBuildingData, building.ID); err != nil {
			return err
		}

		if role == constants.Roles.Owner {
			newBuildingData.Name = building.Name
			newBuildingData.Address = building.Address
			newBuildingData.TotalFloor = building.TotalFloor
			newBuildingData.TotalRoom = len(building.Rooms) + len(building.NewRooms)
			newBuildingData.Rooms = []models.RoomModel{}
			newBuildingData.Images = []models.BuildingImageModel{}
			newBuildingData.Services = []models.BuildingServiceModel{}

			if err := s.buildingRepository.Update(ctx, tx, newBuildingData); err != nil {
				return err
			}

			if len(building.DeletedBuildingImages) > 0 {
				if err := s.buildingRepository.DeleteImages(ctx, tx, building.DeletedBuildingImages); err != nil {
					return err
				}
			}

			if len(building.DeletedSchedules) > 0 {
				if err := s.managerScheduleRepository.Delete(ctx, tx, building.DeletedSchedules); err != nil {
					return err
				}
			}

			if len(building.DeletedRooms) > 0 {
				if err := s.roomRepository.Delete(ctx, tx, building.DeletedRooms); err != nil {
					return err
				}
			}

			if len(building.NewBuildingImages) > 0 {
				newImageNo, err := s.buildingRepository.GetNewImageNo(ctx, building.ID)
				if err != nil {
					return err
				}
				newImage := []models.BuildingImageModel{}
				for index, image := range building.NewBuildingImages {
					filePath, err := utils.StoreFile(image, constants.GetBuildingImageURL("images", buildingIDStr, ""))
					if err != nil {
						return err
					}
					deleteImageList = append(deleteImageList, filePath)
					newImage = append(newImage, models.BuildingImageModel{
						BuildingID: building.ID,
						DefaultFileModel: models.DefaultFileModel{
							Title: filepath.Base(filePath),
							Path:  filePath,
							No:    newImageNo + index,
						},
					})
				}
				if err := s.buildingRepository.AddImage(ctx, tx, &newImage); err != nil {
					return err
				}
			}

			if len(building.NewSchedules) > 0 {
				schedules := []models.ManagerScheduleModel{}
				for _, schedule := range building.NewSchedules {
					startDate, err := utils.ParseTime(schedule.StartDate)
					if err != nil {
						return err
					}

					schedules = append(schedules, models.ManagerScheduleModel{
						BuildingID: building.ID,
						ManagerID:  schedule.ManagerID,
						StartDate:  startDate,
						EndDate:    utils.StringToNullTime(schedule.EndDate),
					})
				}
				if err := s.managerScheduleRepository.Create(ctx, tx, &schedules); err != nil {
					return err
				}
			}

			if len(building.NewRooms) > 0 {
				rooms := []models.RoomModel{}
				for _, room := range building.NewRooms {
					rooms = append(rooms, models.RoomModel{
						No:          room.No,
						Floor:       room.Floor,
						Status:      room.Status,
						Area:        room.Area,
						Description: room.Description,
						BuildingID:  building.ID,
					})
				}
				if err := s.roomRepository.Create(ctx, tx, &rooms); err != nil {
					return err
				}

				roomImages := []models.RoomImageModel{}
				for _, room := range building.NewRooms {
					var targetRoom models.RoomModel
					for _, val := range rooms {
						if val.No == room.No {
							targetRoom = val
						}
					}
					roomNoStr := strconv.Itoa(int(targetRoom.No))

					for index, image := range room.Images {
						filePath, err := utils.StoreFile(image, constants.GetRoomImageURL("images", buildingIDStr, roomNoStr, ""))
						if err != nil {
							return err
						}
						roomImages = append(roomImages, models.RoomImageModel{
							RoomID: targetRoom.ID,
							DefaultFileModel: models.DefaultFileModel{
								Path:  filePath,
								No:    index + 1,
								Title: filepath.Base(filePath),
							},
						})
						deleteImageList = append(deleteImageList, filePath)
					}
				}

				if err := s.roomRepository.CreateImage(ctx, tx, &roomImages); err != nil {
					return err
				}
			}

			{
				scheduleIDs := []int64{}
				for _, schedule := range building.Schedules {
					scheduleIDs = append(scheduleIDs, schedule.ID)
				}
				schedules := []models.ManagerScheduleModel{}
				if err := s.managerScheduleRepository.GetByIDs(ctx, &schedules, scheduleIDs); err != nil {
					return err
				}
				for index, schedule := range schedules {
					for _, val := range building.Schedules {
						if val.ID == schedule.ID {
							startDate, err := utils.ParseTime(val.StartDate)
							if err != nil {
								return err
							}
							schedules[index].ManagerID = val.ManagerID
							schedules[index].StartDate = startDate
							schedules[index].EndDate = utils.StringToNullTime(val.EndDate)
							break
						}
					}
				}
				if err := s.managerScheduleRepository.Update(ctx, tx, &schedules); err != nil {
					return err
				}
			}
		}

		if len(building.DeletedServices) > 0 {
			if err := s.buildingRepository.DeleteServices(ctx, tx, building.DeletedServices); err != nil {
				return err
			}
		}

		if len(building.NewServices) > 0 {
			services := []models.BuildingServiceModel{}
			for _, service := range building.NewServices {
				services = append(services, models.BuildingServiceModel{
					Name:       service.Name,
					Price:      service.Price,
					BuildingID: building.ID,
				})
			}
			if err := s.buildingRepository.AddServices(ctx, tx, &services); err != nil {
				return err
			}
		}

		{
			serviceIDs := []int64{}
			for _, service := range building.Services {
				serviceIDs = append(serviceIDs, service.ID)
			}
			services := []models.BuildingServiceModel{}
			if err := s.buildingRepository.GetServicesByIDs(ctx, &services, serviceIDs); err != nil {
				return err
			}
			for index, service := range services {
				for _, val := range building.Services {
					if val.ID == service.ID {
						services[index].Name = val.Name
						services[index].Price = val.Price
						break
					}
				}
			}
			if err := s.buildingRepository.UpdateServices(ctx, tx, &services); err != nil {
				return err
			}
		}

		if len(building.DeletedRoomImages) > 0 {
			if err := s.roomRepository.DeleteImages(ctx, tx, building.DeletedRoomImages); err != nil {
				return err
			}
		}

		{
			roomIDs := []int64{}
			for _, room := range building.Rooms {
				roomIDs = append(roomIDs, room.ID)
			}
			rooms := []models.RoomModel{}
			if err := s.roomRepository.GetByIDs(ctx, &rooms, roomIDs); err != nil {
				return err
			}
			for index, room := range rooms {
				for _, val := range building.Rooms {
					if val.ID == room.ID {
						rooms[index].No = val.No
						rooms[index].Floor = val.Floor
						rooms[index].Status = val.Status
						rooms[index].Area = val.Area
						rooms[index].Description = val.Description
						rooms[index].Contracts = []models.ContractModel{}
						rooms[index].Images = []models.RoomImageModel{}
						break
					}
				}
			}
			if err := s.roomRepository.Update(ctx, tx, &rooms); err != nil {
				return err
			}

			roomImages := []models.RoomImageModel{}
			for _, room := range building.Rooms {
				var targetRoom models.RoomModel
				for _, val := range rooms {
					if val.ID == room.ID {
						targetRoom = val
						break
					}
				}
				roomNoStr := strconv.Itoa(int(targetRoom.No))

				lastestImageNo, err := s.roomRepository.GetNewImageNo(ctx, targetRoom.ID)
				if err != nil {
					return err
				}

				for index, image := range room.NewImages {
					filePath, err := utils.StoreFile(image, constants.GetRoomImageURL("images", buildingIDStr, roomNoStr, ""))
					if err != nil {
						return err
					}
					roomImages = append(roomImages, models.RoomImageModel{
						RoomID: targetRoom.ID,
						DefaultFileModel: models.DefaultFileModel{
							Path:  filePath,
							No:    lastestImageNo + index,
							Title: filepath.Base(filePath),
						},
					})
					deleteImageList = append(deleteImageList, filePath)
				}
			}

			if len(roomImages) > 0 {
				if err := s.roomRepository.CreateImage(ctx, tx, &roomImages); err != nil {
					return err
				}
			}
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

func (s *BuildingService) CheckManagerPermission(ctx *gin.Context, buildingID int64) bool {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager {
		buildings := []models.BuildingModel{}

		if err := s.buildingRepository.GetBuildingBaseOnSchedule(ctx, &buildings, ctx.GetInt64("userID")); err != nil {
			return false
		}

		if len(buildings) == 0 {
			return false
		}

		var result = false

		for _, building := range buildings {
			if building.ID == buildingID {
				result = true
				break
			}
		}

		return result
	}

	return true
}
