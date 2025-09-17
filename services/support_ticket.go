package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"database/sql"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupportTicketService struct {
	supportTicketRepository *repositories.SupportTicketRepository
	buildingService         *BuildingService
	contractRepository      *repositories.ContractRepository
}

func NewSupportTicketService() *SupportTicketService {
	return &SupportTicketService{
		supportTicketRepository: repositories.NewSupportTicketRepository(),
		buildingService:         NewBuildingService(false),
		contractRepository:      repositories.NewContractRepository(),
	}
}

func (s *SupportTicketService) GetSupportTickets(ctx *gin.Context, tickets *[]structs.SupportTicket, limit, offset int64, startDate string, endDate string) error {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager || role == constants.Roles.Customer {
		userID := ctx.GetInt64("userID")

		if role == constants.Roles.Manager {
			return s.supportTicketRepository.GetSupportTickets(ctx, tickets, limit, offset, startDate, endDate, &userID)
		} else {
			return s.supportTicketRepository.GetTicketsByCustomerID2(ctx, tickets, limit, offset, startDate, endDate, userID)
		}
	}

	return s.supportTicketRepository.GetSupportTickets(ctx, tickets, limit, offset, startDate, endDate, nil)
}

func (s *SupportTicketService) GetSupportTicket(ctx *gin.Context, ticket *structs.SupportTicket, ticketID int64) error {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager || role == constants.Roles.Customer {
		userID := ctx.GetInt64("userID")

		if role == constants.Roles.Manager {
			return s.supportTicketRepository.GetSupportTicket(ctx, ticket, ticketID, &userID)
		} else {
			return s.supportTicketRepository.GetTicketByCustomerID(ctx, ticket, userID, ticketID)
		}
	}

	return s.supportTicketRepository.GetSupportTicket(ctx, ticket, ticketID, nil)
}

func (s *SupportTicketService) CheckManagerPermission(ctx *gin.Context, ticketID int64, managerID int64) bool {
	building := &models.BuildingModel{}
	if err := s.supportTicketRepository.GetTicketBuilding(ctx, ticketID, building); err != nil {
		return false
	}

	return s.buildingService.CheckManagerPermission(ctx, building.ID)
}

func (s *SupportTicketService) ApproveSupportTicket(ctx *gin.Context, ticketID int64) (bool, error) {
	role := ctx.GetString("role")

	userID := ctx.GetInt64("userID")

	ticket := &models.SupportTicketModel{}

	if err := s.supportTicketRepository.GetById(ctx, ticket, ticketID); err != nil {
		return false, err
	}

	if role == constants.Roles.Manager {
		if !s.CheckManagerPermission(ctx, ticketID, userID) {
			return false, nil
		}

		if ticket.ManagerID != 0 {
			return false, nil
		}

		ticket.ManagerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.ManagerID = userID
		ticket.ManagerResult = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
		ticket.Status = constants.Common.SupportTicketStatus.PENDING
	} else {
		if ticket.OwnerID != 0 || (ticket.ManagerID == 0 && !ctx.GetBool("ticketByPass")) {
			return false, nil
		}

		ticket.OwnerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.OwnerID = userID
		ticket.OwnerResult = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
		ticket.Status = constants.Common.SupportTicketStatus.APPROVED
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		if err := s.supportTicketRepository.Update(ctx, tx, ticket, ticketID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *SupportTicketService) DenySupportTicket(ctx *gin.Context, ticketID int64) (bool, error) {
	role := ctx.GetString("role")
	userID := ctx.GetInt64("userID")

	ticket := &models.SupportTicketModel{}

	if err := s.supportTicketRepository.GetById(ctx, ticket, ticketID); err != nil {
		return false, err
	}

	if role == constants.Roles.Manager {
		if !s.CheckManagerPermission(ctx, ticketID, userID) {
			return false, nil
		}

		if ticket.ManagerID != 0 {
			return false, nil
		}

		ticket.ManagerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.ManagerID = userID
		ticket.ManagerResult = sql.NullBool{
			Bool:  false,
			Valid: true,
		}
	} else {
		if ticket.OwnerID != 0 || (ticket.ManagerID == 0 && !ctx.GetBool("ticketByPass")) {
			return false, nil
		}

		ticket.OwnerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.OwnerID = userID
		ticket.OwnerResult = sql.NullBool{
			Bool:  false,
			Valid: true,
		}
	}
	ticket.Status = constants.Common.SupportTicketStatus.REJECTED

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		if err := s.supportTicketRepository.Update(ctx, tx, ticket, ticketID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *SupportTicketService) DeleteTickets(ctx *gin.Context, ids []int64) (bool, error) {
	deletableTickets := &[]models.SupportTicketModel{}
	if err := s.supportTicketRepository.GetDeletableTickets(ctx, deletableTickets, ids, ctx.GetInt64("userID")); err != nil {
		return true, err
	}

	if len(*deletableTickets) != len(ids) {
		return false, nil
	}

	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		if err := s.supportTicketRepository.Delete(ctx, tx, ids); err != nil {
			return err
		}
		return nil
	})
}

func (s *SupportTicketService) UpdateSupportTicket(ctx *gin.Context, ticketID int64, ticket *structs.UpdateSupportTicketRequest) (bool, bool, error) {
	structTicket := &structs.SupportTicket{}
	if err := s.supportTicketRepository.GetTicketByCustomerID(ctx, structTicket, ctx.GetInt64("userID"), ticketID); err != nil {
		return true, true, err
	}

	if structTicket.ID == 0 {
		return true, false, nil
	}

	if !(structTicket.CustomerID == ctx.GetInt64("userID") && structTicket.Status == constants.Common.SupportTicketStatus.PENDING && structTicket.OwnerID == 0 && structTicket.ManagerID == 0) {
		return false, true, nil
	}

	deleteFileList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		if len(ticket.DeletedFiles) > 0 {
			if err := s.supportTicketRepository.DeleteTicketFiles(ctx, tx, ticketID, ticket.DeletedFiles); err != nil {
				return err
			}
		}

		if len(ticket.NewFiles) > 0 {
			ticketIDStr := strconv.Itoa(int(ticketID))
			newTicketFiles := []models.SupportTicketFileModel{}

			for _, file := range ticket.NewFiles {
				filePath, err := utils.StoreFile(file, constants.GetTicketImageURL("images", ticketIDStr, ""))
				if err != nil {
					return err
				}
				newTicketFiles = append(newTicketFiles, models.SupportTicketFileModel{
					SupportTicketID: ticketID,
					DefaultFileModel: models.DefaultFileModel{
						Path:  filePath,
						Title: filepath.Base(filePath),
					},
				})
				deleteFileList = append(deleteFileList, filePath)
			}

			if err := s.supportTicketRepository.AddFile(ctx, tx, &newTicketFiles); err != nil {
				return err
			}
		}

		ticketModel := &models.SupportTicketModel{}
		if err := s.supportTicketRepository.GetById(ctx, ticketModel, ticketID); err != nil {
			return err
		}

		ticketModel.Title = ticket.Title
		ticketModel.Content = ticket.Content

		if err := s.supportTicketRepository.Update(ctx, tx, ticketModel, ticketID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		for _, path := range deleteFileList {
			utils.RemoveFile(path)
		}
		return true, true, err
	}

	return true, true, nil
}

func (s *SupportTicketService) AddSupportTicket(ctx *gin.Context, ticket *structs.CreateSupportTicketRequest) (bool, error) {
	userID := ctx.GetInt64("userID")

	contract := &structs.Contract{}

	if err := s.contractRepository.GetRoomActiveContract(ctx, contract, ticket.RoomID); err != nil {
		return true, err
	}

	if contract.ID == 0 {
		return false, nil
	}

	isAllowed := false

	if contract.HouseholderID == userID {
		isAllowed = true
	}

	for _, resident := range contract.Residents {
		if resident.UserAccountID.Int64 == userID {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return false, nil
	}

	newTicket := &models.SupportTicketModel{
		Title:      ticket.Title,
		Content:    ticket.Content,
		CustomerID: userID,
		ContractID: contract.ID,
	}

	deleteFileList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		if err := s.supportTicketRepository.Add(ctx, tx, newTicket); err != nil {
			return err
		}

		if len(ticket.Files) > 0 {
			ticketIDStr := strconv.Itoa(int(newTicket.ID))
			newTicketFiles := []models.SupportTicketFileModel{}

			for _, file := range ticket.Files {
				filePath, err := utils.StoreFile(file, constants.GetTicketImageURL("images", ticketIDStr, ""))
				if err != nil {
					return err
				}
				newTicketFiles = append(newTicketFiles, models.SupportTicketFileModel{
					SupportTicketID: newTicket.ID,
					DefaultFileModel: models.DefaultFileModel{
						Path:  filePath,
						Title: filepath.Base(filePath),
					},
				})
				deleteFileList = append(deleteFileList, filePath)
			}

			if err := s.supportTicketRepository.AddFile(ctx, tx, &newTicketFiles); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		for _, path := range deleteFileList {
			utils.RemoveFile(path)
		}
		return true, err
	}

	return true, nil
}
