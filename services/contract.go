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

func (s *ContractService) GetContractList(ctx *gin.Context, contracts *[]structs.Contract, limit int64, offset int64) error {
	role, exists := ctx.Get("role")

	if !exists {
		return errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager || role.(string) == constants.Roles.Customer {
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

		if role.(string) == constants.Roles.Manager {
			return s.contractRepository.GetContractsByManagerID2(ctx, contracts, claim.UserID, limit, offset)
		} else {
			return s.contractRepository.GetContractsByCustomerID2(ctx, contracts, claim.UserID, limit, offset)
		}

	}

	return s.contractRepository.GetContracts(ctx, contracts, limit, offset)
}

func (s *ContractService) GetContractDetail(ctx *gin.Context, contract *structs.Contract, id int64) (bool, error) {
	if err := s.contractRepository.GetContractByID(ctx, contract, id); err != nil {
		return false, err
	}

	if contract.ID == 0 {
		return true, nil
	}

	if err := s.billRepository.GetByContractId(ctx, &contract.Bills, contract.ID); err != nil {
		return false, err
	}

	role, exists := ctx.Get("role")

	if !exists {
		return false, errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager || role.(string) == constants.Roles.Customer {
		jwt, exists := ctx.Get("jwt")

		if !exists {
			return false, errors.New("jwt not found")
		}

		token, err := utils.ValidateJWTToken(jwt.(string))

		if err != nil {
			return false, err
		}

		claim := &structs.JTWClaim{}

		utils.ExtractJWTClaim(token, claim)

		if role.(string) == constants.Roles.Manager {
			// return contract.CreatorID == claim.UserID, nil
			return true, nil
		} else {
			return contract.HouseholderID == claim.UserID, nil
		}
	}

	return true, nil
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

func (s *ContractService) DeleteContract2(ctx *gin.Context, IDs []int64) (bool, error) {
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
		if err := s.contractRepository.GetDeletableContracts2(ctx, &contracts, IDs, &claim.UserID); err != nil {
			return true, err
		}

		if len(contracts) != len(IDs) {
			return false, nil
		}
	} else if role.(string) == constants.Roles.Owner {
		contracts := []models.ContractModel{}
		if err := s.contractRepository.GetDeletableContracts2(ctx, &contracts, IDs, nil); err != nil {
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
