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

type RoomService struct {
	contractService         *ContractService
	roomRepository          *repositories.RoomRepository
	contractRepository      *repositories.ContractRepository
	supportTicketRepository *repositories.SupportTicketRepository
}

func NewRoomService() *RoomService {
	return &RoomService{
		contractService:         NewContractService(),
		roomRepository:          repositories.NewRoomRepository(),
		contractRepository:      repositories.NewContractRepository(),
		supportTicketRepository: repositories.NewSupportTicketRepository(),
	}
}

// func (s *RoomService) DeleteWithoutTransaction(ctx *gin.Context, tx *gorm.DB, id []int64) error {
// 	contractIDs := []int64{}
// 	contracts := []models.ContractModel{}
// 	if err := s.contractRepository.GetContractByRoomID(ctx, &contracts, id); err != nil {
// 		return err
// 	}

// 	for _, contract := range contracts {
// 		contractIDs = append(contractIDs, contract.ID)
// 	}

// 	if err := s.roomRepository.Delete(ctx, tx, id); err != nil {
// 		return err
// 	}

// 	if err := s.contractService.DeleteWithoutTransaction(ctx, tx, contractIDs); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *RoomService) GetRoomDetail(ctx *gin.Context, room *models.RoomModel, id int64) error {
	if err := s.roomRepository.GetById(ctx, room, id); err != nil {
		return err
	}
	return nil
}

func (s *RoomService) GetRoomByRoomIDAndBuildingID(ctx *gin.Context, room *models.RoomModel, roomID int64, buildingID int64) error {
	if err := s.roomRepository.GetRoomByRoomIDAndBuildingID(ctx, room, roomID, buildingID); err != nil {
		return err
	}
	return nil
}

func (s *RoomService) GetContractByRoomIDAndBuildingID(ctx *gin.Context, contracts *[]structs.Contract, roomID int64, buildingID int64) error {
	// role, exists := ctx.Get("role")

	// if !exists {
	// 	return errors.New("role not found")
	// }

	// if role.(string) == constants.Roles.Manager {
	// 	jwt, exists := ctx.Get("jwt")

	// 	if !exists {
	// 		return errors.New("jwt not found")
	// 	}

	// 	token, err := utils.ValidateJWTToken(jwt.(string))

	// 	if err != nil {
	// 		return err
	// 	}

	// 	claim := &structs.JTWClaim{}

	// 	utils.ExtractJWTClaim(token, claim)

	// 	return s.contractRepository.GetContractByRoomIDAndBuildingIDAndManagerID(ctx, contracts, roomID, buildingID, claim.UserID)
	// }

	return s.contractRepository.GetContractByRoomIDAndBuildingID(ctx, contracts, roomID, buildingID)
}

func (s *RoomService) GetTicketByRoomIDAndBuildingID(ctx *gin.Context, roomID int64, buildingID int64, startDate string, endDate string, tickets *[]models.SupportTicketModel) error {
	if err := s.supportTicketRepository.GetTicketByRoomIDAndBuildingID(ctx, roomID, buildingID, startDate, endDate, tickets); err != nil {
		return err
	}
	return nil
}

func (s *RoomService) UpdateRoomByRoomIDAndBuildingID(ctx *gin.Context, oldRoomData *models.RoomModel, room *structs.EditRoom2, roomID int64, buildingID int64) error {
	deleteImageList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		buildingIDStr := strconv.Itoa(int(buildingID))
		roomNoStr := strconv.Itoa(int(oldRoomData.No))
		roomImages := []models.RoomImageModel{}

		lastestImageNo, err := s.roomRepository.GetNewImageNo(ctx, oldRoomData.ID)
		if err != nil {
			return err
		}

		if len(room.DeletedRoomImages) > 0 {
			if err := s.roomRepository.DeleteImages(ctx, tx, room.DeletedRoomImages); err != nil {
				return err
			}
		}

		for index, image := range room.NewRoomImages {
			filePath, err := utils.StoreFile(image, constants.GetRoomImageURL("images", buildingIDStr, roomNoStr, ""))
			if err != nil {
				return err
			}
			roomImages = append(roomImages, models.RoomImageModel{
				RoomID: oldRoomData.ID,
				DefaultFileModel: models.DefaultFileModel{
					Path:  filePath,
					No:    lastestImageNo + index,
					Title: filepath.Base(filePath),
				},
			})
			deleteImageList = append(deleteImageList, filePath)
		}

		if len(roomImages) > 0 {
			if err := s.roomRepository.CreateImage(ctx, tx, &roomImages); err != nil {
				return err
			}
		}

		oldRoomData.Status = room.Status
		oldRoomData.Description = room.Description
		oldRoomData.Area = room.Area

		newRoomData := []models.RoomModel{}

		newRoomData = append(newRoomData, *oldRoomData)

		if err := s.roomRepository.Update(ctx, tx, &newRoomData); err != nil {
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

func (s *RoomService) UpdateRoomStatus() error {
	return config.WorkerDB.Transaction(func(tx *gorm.DB) error {
		if err := s.roomRepository.UpdateRoomStatus(tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *RoomService) GetRoomList(ctx *gin.Context, rooms *[]structs.BuildingRoom) error {
	jwt, exists := ctx.Get("jwt")

	if !exists {
		return errors.New("jwt not found")
	}

	token, err := utils.ValidateJWTToken(jwt.(string))

	if err != nil {
		return err
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	if err := s.roomRepository.GetRoomList(ctx, rooms, claim.UserID); err != nil {
		return err
	}

	return nil
}
