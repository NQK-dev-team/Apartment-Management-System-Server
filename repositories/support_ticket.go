package repositories

import (
	"api/config"
	"api/models"
	"api/structs"
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

func (r *SupportTicketRepository) GetSupportTickets(ctx *gin.Context, tickets *[]structs.SupportTicket, limit int64, offset int64, startDate string, endDate string, isOwner bool, managerID *int64) error {
	if isOwner {
		if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").
			Where("created_at::timestamp::date >= ? AND created_at::timestamp::date <= ? AND manager_id IS NOT NULL", startDate, endDate).
			Limit(int(limit)).Offset(int(offset)).Order("created_at desc, owner_resolve_time desc, manager_resolve_time desc").
			Find(tickets).Error; err != nil {
			return err
		}
	} else {
		if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").
			Where("created_at::timestamp::date >= ? AND created_at::timestamp::date <= ? AND (manager_id IS NOT NULL OR manager_id = ?)", startDate, endDate, *managerID).
			Limit(int(limit)).Offset(int(offset)).Order("created_at desc, owner_resolve_time desc, manager_resolve_time desc").
			Find(tickets).Error; err != nil {
			return err
		}
	}

	for i := range *tickets {
		if err := config.DB.Raw("SELECT room.no AS room_no FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].RoomNo).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].RoomFloor).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT building.name AS building_name FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].BuildingName).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketsByManagerID(ctx *gin.Context, tickets *[]structs.SupportTicket, managerID int64, limit int64, offset int64, startDate string, endDate string) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").
		// Joins("JOIN support_ticket ON support_ticket.id = manager_resolve_support_ticket.support_ticket_id").
		Where("manager_id = ? AND created_at::timestamp::date >= ? AND created_at::timestamp::date <= ?", managerID, startDate, endDate).
		Limit(int(limit)).Offset(int(offset)).Order("created_at desc, owner_resolve_time desc, manager_resolve_time desc").
		Find(tickets).Error; err != nil {
		return err
	}

	for i := range *tickets {
		// if err := config.DB.Raw("SELECT id FROM \"user\" JOIN manager_resolve_support_ticket ON manager_resolve_support_ticket.manager_id=\"user\".id WHERE NOT manager_id = ? AND support_ticket_id = ?", managerID, (*tickets)[i].SupportTicket.ID).Scan(&(*tickets)[i].OwnerID).Error; err != nil {
		// 	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 		return err
		// 	}
		// } else if (*tickets)[i].OwnerID != 0 {
		// 	if err := config.DB.Model(&models.UserModel{}).Where("id = ?", (*tickets)[i].OwnerID).First(&(*tickets)[i].Owner).Error; err != nil {
		// 		return err
		// 	}

		// 	if err := config.DB.Raw("SELECT result FROM manager_resolve_support_ticket WHERE manager_id = ? AND support_ticket_id = ?", (*tickets)[i].OwnerID, (*tickets)[i].SupportTicket.ID).Scan(&(*tickets)[i].OwnerResult).Error; err != nil {
		// 		return err
		// 	}

		// 	if err := config.DB.Raw("SELECT resolve_time FROM manager_resolve_support_ticket WHERE manager_id = ? AND support_ticket_id = ?", (*tickets)[i].OwnerID, (*tickets)[i].SupportTicket.ID).Scan(&(*tickets)[i].OwnerResolveTime).Error; err != nil {
		// 		return err
		// 	}
		// }

		// if err := config.DB.Model(&models.SupportTicketFileModel{}).Where("support_ticket_id = ?", (*tickets)[i].SupportTicket.ID).First(&(*tickets)[i].SupportTicket.Files).Error; err != nil {
		// 	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 		return err
		// 	}
		// }

		// if err := config.DB.Model(&models.UserModel{}).Where("id = ?", (*tickets)[i].SupportTicket.CustomerID).First(&(*tickets)[i].SupportTicket.Customer).Error; err != nil {
		// 	return err
		// }

		if err := config.DB.Raw("SELECT room.no AS room_no FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].RoomNo).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].RoomFloor).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT building.name AS building_name FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].BuildingName).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketsByCustomerID(ctx *gin.Context, tickets *[]structs.SupportTicket, customerID int64) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").
		Where("customer_id = ?", customerID).
		Order("created_at desc, owner_resolve_time desc, manager_resolve_time desc").
		Find(tickets).Error; err != nil {
		return err
	}

	for i := range *tickets {
		if err := config.DB.Raw("SELECT room.no AS room_no FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].RoomNo).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].RoomFloor).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT building.name AS building_name FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*tickets)[i].ID).Scan(&(*tickets)[i].BuildingName).Error; err != nil {
			return err
		}
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

func (r *SupportTicketRepository) Update(ctx *gin.Context, tx *gorm.DB, ticket *models.SupportTicketModel, id int64) error {
	userID := ctx.GetInt64("userID")
	if ticket.OwnerID != 0 {
		if err := tx.Set("userID", userID).Model(&models.SupportTicketModel{}).Where("id = ?", id).Updates(ticket).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Set("userID", userID).Model(&models.SupportTicketModel{}).Omit("OwnerID").Where("id = ?", id).Updates(ticket).Error; err != nil {
			return err
		}
	}
	return nil
}
