package services

import (
	"api/models"
	"api/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomService struct {
	contractService    *ContractService
	roomRepository     *repositories.RoomRepository
	contractRepository *repositories.ContractRepository
}

func NewRoomService() *RoomService {
	return &RoomService{
		contractService:    NewContractService(),
		roomRepository:     repositories.NewRoomRepository(),
		contractRepository: repositories.NewContractRepository(),
	}
}

func (s *RoomService) DeleteWithoutTransaction(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	contractIDs := []int64{}
	contracts := []models.ContractModel{}
	if err := s.contractRepository.GetContractByRoomID(ctx, &contracts, id); err != nil {
		return err
	}

	for _, contract := range contracts {
		contractIDs = append(contractIDs, contract.ID)
	}

	if err := s.roomRepository.Delete(ctx, tx, id); err != nil {
		return err
	}

	if err := s.contractService.DeleteWithoutTransaction(ctx, tx, contractIDs); err != nil {
		return err
	}

	return nil
}

func (s *RoomService) GetRoomDetail(ctx *gin.Context, room *models.RoomModel, id int64) error {
	if err := s.roomRepository.GetById(ctx, room, id); err != nil {
		return err
	}
	return nil
}
