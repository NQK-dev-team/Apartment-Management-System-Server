package services

import (
	"api/repositories"
	"api/structs"

	"github.com/gin-gonic/gin"
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

func (s *ContractService) GetContractByRoomIDAndBuildingID(ctx *gin.Context, contracts *[]structs.Contract, roomID int64, buildingID int64) error {
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
