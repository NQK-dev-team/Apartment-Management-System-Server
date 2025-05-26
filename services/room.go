package services

import (
	"api/models"
	"api/repositories"
	"api/structs"

	"github.com/gin-gonic/gin"
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
