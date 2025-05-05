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
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupportTicketService struct {
	supportTicketRepository *repositories.SupportTicketRepository
}

func NewSupportTicketService() *SupportTicketService {
	return &SupportTicketService{
		supportTicketRepository: repositories.NewSupportTicketRepository(),
	}
}

func (s *SupportTicketService) GetSupportTickets(ctx *gin.Context, tickets *[]structs.SupportTicket, limit, offset int64, startDate string, endDate string) error {
	if err := s.supportTicketRepository.GetSupportTickets(ctx, tickets, limit, offset, startDate, endDate); err != nil {
		return err
	}
	return nil
}

func (s *SupportTicketService) ApproveSupportTicket(ctx *gin.Context, ticketID int64) (bool, error) {
	role, exists := ctx.Get("role")

	if !exists {
		return false, errors.New("role not found")
	}

	jwt, exists := ctx.Get("jwt")

	if !exists {
		return false, errors.New("jwt not found")
	}

	token, err := utils.ValidateJWTToken(jwt.(string))

	if err != nil {
		return false, errors.New("jwt not valid")
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	ticket := &models.SupportTicketModel{}

	if err := s.supportTicketRepository.GetById(ctx, ticket, ticketID); err != nil {
		return false, err
	}

	if role.(string) == constants.Roles.Manager {
		if ticket.ManagerID != 0 {
			return false, nil
		}

		ticket.ManagerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.ManagerID = claim.UserID
		ticket.ManagerResult = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
		ticket.Status = 1
	} else {
		if ticket.OwnerID != 0 {
			return false, nil
		}

		ticket.OwnerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.OwnerID = claim.UserID
		ticket.OwnerResult = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
		ticket.Status = 2
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
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
	role, exists := ctx.Get("role")

	if !exists {
		return false, errors.New("role not found")
	}

	jwt, exists := ctx.Get("jwt")

	if !exists {
		return false, errors.New("jwt not found")
	}

	token, err := utils.ValidateJWTToken(jwt.(string))

	if err != nil {
		return false, errors.New("jwt not valid")
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	ticket := &models.SupportTicketModel{}

	if err := s.supportTicketRepository.GetById(ctx, ticket, ticketID); err != nil {
		return false, err
	}

	if role.(string) == constants.Roles.Manager {
		if ticket.ManagerID != 0 {
			return false, nil
		}

		ticket.ManagerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.ManagerID = claim.UserID
		ticket.ManagerResult = sql.NullBool{
			Bool:  false,
			Valid: true,
		}
	} else {
		if ticket.OwnerID != 0 {
			return false, nil
		}

		ticket.OwnerResolveTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		ticket.OwnerID = claim.UserID
		ticket.OwnerResult = sql.NullBool{
			Bool:  false,
			Valid: true,
		}
	}
	ticket.Status = 3

	err = config.DB.Transaction(func(tx *gorm.DB) error {
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
