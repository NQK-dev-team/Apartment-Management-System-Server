package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractService struct {
	contractRepository *repositories.ContractRepository
	billRepository     *repositories.BillRepository
	ticketRepository   *repositories.SupportTicketRepository
}

func NewContractService() *ContractService {
	return &ContractService{
		contractRepository: repositories.NewContractRepository(),
		billRepository:     repositories.NewBillRepository(),
		ticketRepository:   repositories.NewSupportTicketRepository(),
	}
}

func (s *ContractService) DeleteContract(ctx *gin.Context, IDs []int64, roomID int64, buildingID int64) (bool, error) {
	role, exists := ctx.Get("role")

	if !exists {
		return true, errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager {
		jwt, exists := ctx.Get("jwt")

		if !exists {
			return true, errors.New("jwt not found")
		}

		token, err := utils.ValidateJWTToken(jwt.(string))

		if err != nil {
			return true, err
		}

		claim := &structs.JTWClaim{}

		utils.ExtractJWTClaim(token, claim)

		contracts := []models.ContractModel{}
		if err := s.contractRepository.GetDeletableContracts(ctx, &contracts, IDs, &claim.UserID, roomID, buildingID); err != nil {
			return true, err
		}

		if len(contracts) != len(IDs) {
			return false, nil
		}
	} else if role.(string) == constants.Roles.Owner {
		contracts := []models.ContractModel{}
		if err := s.contractRepository.GetDeletableContracts(ctx, &contracts, IDs, nil, roomID, buildingID); err != nil {
			return true, err
		}

		if len(contracts) != len(IDs) {
			return false, nil
		}
	}
	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		return s.contractRepository.Delete(ctx, tx, IDs)
	})
}

// func (s *ContractService) DeleteWithoutTransaction(ctx *gin.Context, tx *gorm.DB, id []int64) error {
// 	contracts := []models.ContractModel{}
// 	if err := s.contractRepository.GetContractByIDs(ctx, &contracts, id); err != nil {
// 		return err
// 	}

// 	billIDs := []int64{}
// 	ticketIDs := []int64{}

// 	for _, contract := range contracts {
// 		for _, bill := range contract.Bills {
// 			billIDs = append(billIDs, bill.ID)
// 		}

// 		for _, ticket := range contract.SupportTickets {
// 			ticketIDs = append(ticketIDs, ticket.ID)
// 		}
// 	}

// 	if err := s.contractRepository.Delete(ctx, tx, id); err != nil {
// 		return err
// 	}

// 	if err := s.billRepository.Delete(ctx, tx, billIDs); err != nil {
// 		return err
// 	}

// 	if err := s.ticketRepository.Delete(ctx, tx, ticketIDs); err != nil {
// 		return err
// 	}

// 	return nil
// }
