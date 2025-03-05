package services

import (
	"api/repositories"
)

type SupportTicketService struct {
	supportTicketRepository *repositories.SupportTicketRepository
}

func NewSupportTicketService() *SupportTicketService {
	return &SupportTicketService{
		supportTicketRepository: repositories.NewSupportTicketRepository(),
	}
}
