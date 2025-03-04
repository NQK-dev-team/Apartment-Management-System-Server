package services

import (
	"api/models"
	"api/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type ContractService struct {
	billService          *BillService
	supportTicketService *SupportTicketService
	contractRepository   *repositories.ContractRepository
}

func NewContractService() *ContractService {
	return &ContractService{
		billService:          NewBillService(),
		supportTicketService: NewSupportTicketService(),
		contractRepository:   repositories.NewContractRepository(),
	}
}

func (s *ContractService) DeleteWithoutTransaction(ctx *gin.Context, id int64) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}

	now := time.Now()

	contract := &models.ContractModel{}
	if err := s.contractRepository.GetById(ctx, contract, id); err != nil {
		return err
	}

	contract.DefaultModel.DeletedBy = userID.(int64)
	contract.DefaultModel.DeletedAt.Valid = true
	contract.DefaultModel.DeletedAt.Time = now

	for index := range contract.Files {
		contract.Files[index].DefaultFileModel.DeletedBy = userID.(int64)
		contract.Files[index].DefaultFileModel.DeletedAt.Valid = true
		contract.Files[index].DefaultFileModel.DeletedAt.Time = now
	}

	if err := s.contractRepository.QuietUpdate(ctx, contract); err != nil {
		return err
	}

	for _, bill := range contract.Bills {
		if err := s.billService.DeleteWithoutTransaction(ctx, bill.ID); err != nil {
			return err
		}
	}

	for _, ticket := range contract.SupportTickets {
		if err := s.supportTicketService.DeleteWithoutTransaction(ctx, ticket.ID); err != nil {
			return err
		}
	}

	return nil
}
