package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractService struct {
	contractRepository *repositories.ContractRepository
	billRepository     *repositories.BillRepository
	ticketRepository   *repositories.SupportTicketRepository
	buildingRepository *repositories.BuildingRepository
	userRepository     *repositories.UserRepository
}

func NewContractService() *ContractService {
	return &ContractService{
		contractRepository: repositories.NewContractRepository(),
		billRepository:     repositories.NewBillRepository(),
		ticketRepository:   repositories.NewSupportTicketRepository(),
		buildingRepository: repositories.NewBuildingRepository(),
		userRepository:     repositories.NewUserRepository(),
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
			isAllowed, err := s.CheckManagerContractPermission(ctx, claim.UserID, contract.ID)

			if err != nil {
				return false, err
			}

			if !isAllowed {
				return false, nil
			}

			// return contract.CreatorID == claim.UserID, nil
			return true, nil
		} else {
			return contract.HouseholderID == claim.UserID, nil
		}
	}

	return true, nil
}

func (s *ContractService) GetContractBill(ctx *gin.Context, bills *[]models.BillModel, contractID int64) (bool, error) {
	contract := &structs.Contract{}

	if err := s.contractRepository.GetContractByID(ctx, contract, contractID); err != nil {
		return false, err
	}

	if contractID == 0 {
		return true, nil
	}

	if err := s.billRepository.GetByContractId(ctx, bills, contractID); err != nil {
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
			isAllowed, err := s.CheckManagerContractPermission(ctx, claim.UserID, contractID)

			if err != nil {
				return false, err
			}

			if !isAllowed {
				return false, nil
			}

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

func (s *ContractService) CheckManagerContractPermission(ctx *gin.Context, managerID int64, contractID int64) (bool, error) {
	managerBuildings := []models.BuildingModel{}
	if err := s.buildingRepository.GetBuildingBaseOnSchedule(ctx, &managerBuildings, managerID); err != nil {
		return false, err
	}

	contractBuilding := models.BuildingModel{}
	if err := s.buildingRepository.GetBuildingByContractID(ctx, &contractBuilding, contractID); err != nil {
		return false, err
	}

	for _, building := range managerBuildings {
		if building.ID == contractBuilding.ID {
			return true, nil
		}
	}

	return false, nil
}

func (s *ContractService) UpdateContract(ctx *gin.Context, contract *structs.EditContract, contractID int64) (bool, bool, error) {
	oldContractData := &structs.Contract{}

	if err := s.contractRepository.GetContractByID(ctx, oldContractData, contractID); err != nil {
		return true, true, err
	}

	if oldContractData.ID == 0 {
		return true, true, errors.New("contract not found")
	}

	role, exists := ctx.Get("role")

	if !exists {
		return true, true, errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager {
		jwt, exists := ctx.Get("jwt")

		if !exists {
			return true, true, errors.New("jwt not found")
		}

		token, err := utils.ValidateJWTToken(jwt.(string))

		if err != nil {
			return true, true, err
		}

		claim := &structs.JTWClaim{}

		utils.ExtractJWTClaim(token, claim)

		isAllowed, err := s.CheckManagerContractPermission(ctx, claim.UserID, contractID)

		if err != nil {
			return true, true, err
		}

		if !isAllowed {
			return false, true, nil
		}
	}

	if oldContractData.Status == constants.Common.ContractStatus.EXPIRED || oldContractData.Status == constants.Common.ContractStatus.CANCELLED {
		return false, true, nil
	}

	isContractSigned := false
	contractActiveResult, err := utils.CompareDates(oldContractData.StartDate.Format("2006-01-02"), time.Now().Format("2006-01-02"))

	if err != nil {
		return true, true, err
	}
	isContractActive := contractActiveResult == 0 || contractActiveResult == -1

	if oldContractData.SignDate.Valid {
		if contract.NewSignDate != "" {
			return true, false, nil
		}
		isContractSigned = true
	} else {
		if contract.NewSignDate != "" {
			result, err := utils.CompareDates(contract.NewSignDate, oldContractData.StartDate.Format("2006-01-02"))
			if err != nil {
				return true, true, err
			}

			if result == 1 {
				return true, false, nil
			}

			isContractSigned = true
		}
	}

	if isContractSigned && isContractActive {
		if contract.Status != constants.Common.ContractStatus.ACTIVE && contract.Status != constants.Common.ContractStatus.EXPIRED && contract.Status != constants.Common.ContractStatus.CANCELLED {
			return true, false, nil
		}

		if oldContractData.Type == constants.Common.ContractType.BUY && contract.Status == constants.Common.ContractStatus.EXPIRED {
			return true, false, nil
		}
	} else if isContractSigned && !isContractActive {
		if contract.Status != constants.Common.ContractStatus.NOT_IN_EFFECT && contract.Status != constants.Common.ContractStatus.CANCELLED {
			return true, false, nil
		}
	} else if !isContractSigned && isContractActive {
		if contract.Status != constants.Common.ContractStatus.CANCELLED {
			return true, false, nil
		}
	} else {
		if contract.Status != constants.Common.ContractStatus.WAITING_FOR_SIGNATURE && contract.Status != constants.Common.ContractStatus.CANCELLED {
			return true, false, nil
		}
	}

	contractIDStr := strconv.Itoa(int(contractID))
	deleteFileList := []string{}

	updatedContractData := &models.ContractModel{}

	if err := s.contractRepository.GetById(ctx, updatedContractData, contractID); err != nil {
		return true, true, err
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		newContractFiles := []models.ContractFileModel{}

		for index, file := range contract.NewFiles {
			filePath, err := utils.StoreFile(file.File, constants.GetContractFileURL("files", contractIDStr, ""))
			if err != nil {
				return err
			}
			newContractFiles = append(newContractFiles, models.ContractFileModel{
				ContractID: contractID,
				DefaultFileModel: models.DefaultFileModel{
					Path:  filePath,
					No:    index + 1,
					Title: file.Title,
				},
			})
			deleteFileList = append(deleteFileList, filePath)
		}

		if len(newContractFiles) > 0 {
			if err := s.contractRepository.AddFile(ctx, tx, &newContractFiles); err != nil {
				return err
			}
		}

		if len(contract.RemovedResidents) > 0 {
			if err := s.contractRepository.DeleteResident(ctx, tx, contract.RemovedResidents); err != nil {
				return err
			}
		}

		newResidents := []models.RoomResidentModel{}
		updateResidents := []models.RoomResidentModel{}

		for _, resident := range contract.Residents {
			var residentData *models.RoomResidentModel

			if resident.UserAccountID != 0 {
				customer := &models.UserModel{}
				if err := s.userRepository.GetCustomerDetail(ctx, customer, resident.UserAccountID); err != nil {
					return err
				}

				if customer.ID == 0 {
					return errors.New("customer not found")
				}

				residentData = &models.RoomResidentModel{
					FirstName:               customer.FirstName,
					LastName:                customer.LastName,
					MiddleName:              customer.MiddleName,
					SSN:                     customer.SSN,
					OldSSN:                  customer.OldSSN,
					DOB:                     customer.DOB,
					POB:                     customer.POB,
					Phone:                   sql.NullString{String: customer.Phone, Valid: customer.Phone != ""},
					Email:                   sql.NullString{String: customer.Email, Valid: customer.Email != ""},
					Gender:                  customer.Gender,
					RelationWithHouseholder: resident.RelationWithHouseholder,
					UserAccountID:           sql.NullInt64{Int64: customer.ID, Valid: true},
				}
			} else {
				residentData = &models.RoomResidentModel{
					FirstName:               resident.FirstName,
					LastName:                resident.LastName,
					MiddleName:              sql.NullString{String: resident.MiddleName, Valid: resident.MiddleName != ""},
					SSN:                     resident.SSN,
					OldSSN:                  sql.NullString{String: resident.OldSSN, Valid: resident.OldSSN != ""},
					DOB:                     resident.DOB,
					POB:                     resident.POB,
					Phone:                   sql.NullString{String: resident.Phone, Valid: resident.Phone != ""},
					Email:                   sql.NullString{String: resident.Email, Valid: resident.Email != ""},
					Gender:                  resident.Gender,
					RelationWithHouseholder: resident.RelationWithHouseholder,
					UserAccountID:           sql.NullInt64{Int64: 0, Valid: false},
				}
			}

			if resident.ID != 0 {
				residentData.ID = resident.ID
				updateResidents = append(updateResidents, *residentData)
			} else {
				newResidents = append(newResidents, *residentData)
			}
		}

		if len(newResidents) > 0 {
			for _, newResident := range newResidents {
				if err := s.contractRepository.AddNewRoomResident(ctx, tx, &newResident, contractID); err != nil {
					return err
				}
			}
		}

		if len(updateResidents) > 0 {
			for _, updateResident := range updateResidents {
				if err := s.contractRepository.UpdateRoomResident(ctx, tx, &updateResident); err != nil {
					return err
				}
			}
		}

		updatedContractData.Status = contract.Status
		if contract.NewSignDate != "" {
			updatedContractData.SignDate = utils.StringToNullTime(contract.NewSignDate)
		}

		if contract.Status == constants.Common.ContractStatus.EXPIRED || contract.Status == constants.Common.ContractStatus.CANCELLED {
			updatedContractData.EndDate = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		}

		if err := s.contractRepository.UpdateContract(ctx, tx, updatedContractData); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		for _, path := range deleteFileList {
			utils.RemoveFile(path)
		}
		return false, false, err
	}

	return true, true, nil
}

func (s *ContractService) UpdateContractStatus() error {
	return config.DBNoLog.Transaction(func(tx *gorm.DB) error {
		return s.contractRepository.UpdateContractStatus(tx)
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
