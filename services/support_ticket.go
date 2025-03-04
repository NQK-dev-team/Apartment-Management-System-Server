package services

import (
	"api/models"
	"api/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type SupportTicketService struct {
	supportTicketRepository *repositories.SupportTicketRepository
}

func NewSupportTicketService() *SupportTicketService {
	return &SupportTicketService{
		supportTicketRepository: repositories.NewSupportTicketRepository(),
	}
}

func (s *SupportTicketService) DeleteWithoutTransaction(ctx *gin.Context, id int64) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}

	now := time.Now()

	ticket := &models.SupportTicketModel{}
	if err := s.supportTicketRepository.GetById(ctx, ticket, id); err != nil {
		return err
	}

	ticket.DefaultModel.DeletedBy = userID.(int64)
	ticket.DefaultModel.DeletedAt.Valid = true
	ticket.DefaultModel.DeletedAt.Time = now

	for index := range ticket.Files {
		ticket.Files[index].DefaultFileModel.DeletedBy = userID.(int64)
		ticket.Files[index].DefaultFileModel.DeletedAt.Valid = true
		ticket.Files[index].DefaultFileModel.DeletedAt.Time = now
	}

	if err := s.supportTicketRepository.QuietUpdate(ctx, ticket); err != nil {
		return err
	}

	return nil
}
