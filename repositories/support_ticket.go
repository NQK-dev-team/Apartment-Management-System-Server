package repositories

import (
	"api/config"
	"api/models"
	"errors"
	"time"

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

func (r *SupportTicketRepository) GetTicketsByManagerID(ctx *gin.Context, tickets *[]models.ManagerResolveSupportTicketModel, managerID int64, limit int64, offset int64) error {
	if err := config.DB.Model(&models.ManagerResolveSupportTicketModel{}).Preload("SupportTicket").Preload("SupportTicket.Files").Preload("SupportTicket.Customer").
		Joins("JOIN support_ticket ON support_ticket.id = manager_resolve_support_ticket.support_ticket_id").
		Where("manager_id = ?", managerID).Limit(int(limit)).Offset(int(offset)).Order("support_ticket.created_at desc, resolve_time desc").
		Find(tickets).Error; err != nil {
		return err
	}
	return nil
}

func (r *SupportTicketRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.SupportTicketModel{}).Where("id IN ?", id).UpdateColumns(models.SupportTicketModel{
		DefaultModel: models.DefaultModel{
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
			DeletedBy: userID,
		},
	}).Error; err != nil {
		return err
	}

	if err := tx.Set("isQuiet", true).Model(&models.SupportTicketFileModel{}).Where("support_ticket_id IN ?", id).UpdateColumns(models.SupportTicketFileModel{
		DefaultFileModel: models.DefaultFileModel{
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
			DeletedBy: userID,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}
