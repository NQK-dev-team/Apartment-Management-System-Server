package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupportTicketRepository struct {
}

func NewSupportTicketRepository() *SupportTicketRepository {
	return &SupportTicketRepository{}
}

func (r *SupportTicketRepository) GetById(ctx *gin.Context, ticket *models.SupportTicketModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("Files").First(ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *SupportTicketRepository) QuietUpdate(ctx *gin.Context, ticket *models.SupportTicketModel) error {
	if err := config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(ticket).Error; err != nil {
		return err
	}
	return nil
}
