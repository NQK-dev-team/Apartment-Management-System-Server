package services

import (
	"api/models"
	"api/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type RoomService struct {
	contractService *ContractService
	roomRepository  *repositories.RoomRepository
}

func NewRoomService() *RoomService {
	return &RoomService{
		contractService: NewContractService(),
		roomRepository:  repositories.NewRoomRepository(),
	}
}

func (s *RoomService) DeleteWithoutTransaction(ctx *gin.Context, id int64) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}

	room := &models.RoomModel{}
	if err := s.roomRepository.GetById(ctx, room, id); err != nil {
		return err
	}

	now := time.Now()

	room.DefaultModel.DeletedBy = userID.(int64)
	room.DefaultModel.DeletedAt.Valid = true
	room.DefaultModel.DeletedAt.Time = now

	for index := range room.Images {
		room.Images[index].DefaultFileModel.DeletedBy = userID.(int64)
		room.Images[index].DefaultFileModel.DeletedAt.Valid = true
		room.Images[index].DefaultFileModel.DeletedAt.Time = now
	}

	if err := s.roomRepository.QuietUpdate(ctx, room); err != nil {
		return err
	}

	for _, contract := range room.Contracts {
		if err := s.contractService.DeleteWithoutTransaction(ctx, contract.ID); err != nil {
			return err
		}
	}

	return nil
}
