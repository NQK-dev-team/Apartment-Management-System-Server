package repositories

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/structs"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupportTicketRepository struct {
}

func NewSupportTicketRepository() *SupportTicketRepository {
	return &SupportTicketRepository{}
}

func (r *SupportTicketRepository) GetTicketBuilding(ctx *gin.Context, ticketID int64, building *models.BuildingModel) error {
	if err := config.DB.Model(&models.BuildingModel{}).
		Joins("JOIN room ON room.building_id = building.id AND room.deleted_at IS NULL").
		Joins("JOIN contract ON contract.room_id = room.id AND contract.deleted_at IS NULL").
		Joins("JOIN support_ticket ON support_ticket.contract_id = contract.id AND support_ticket.deleted_at IS NULL").
		Where("support_ticket.id = ? AND building.deleted_at IS NULL", ticketID).
		Find(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *SupportTicketRepository) GetById(ctx *gin.Context, ticket *models.SupportTicketModel, id int64) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).
		Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
		Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
		Where("support_ticket.id = ? AND support_ticket.deleted_at IS NULL", id).Preload("Files").Find(ticket).Error; err != nil {
		return err
	}
	return nil
}

func (r *SupportTicketRepository) GetSupportTickets(ctx *gin.Context, tickets *[]structs.SupportTicket, limit int64, offset int64, startDate string, endDate string, managerID *int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Distinct().Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
			Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
			Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
			Where("support_ticket.created_at::timestamp::date >= ? AND support_ticket.created_at::timestamp::date <= ? AND support_ticket.manager_id IS NOT NULL", startDate, endDate).
			Limit(int(limit)).Offset(int(offset)).Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
			Find(tickets).Error; err != nil {
			return err
		}
	} else {
		if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Distinct().Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
			Joins("JOIN contract ON contract.id = support_ticket.contract_id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("support_ticket.created_at::timestamp::date >= ? AND support_ticket.created_at::timestamp::date <= ? AND (manager_id = ? OR (manager_id IS NULL AND building.id IN (?))) AND support_ticket.deleted_at IS NULL", startDate, endDate, *managerID,
				config.DB.Model(&models.BuildingModel{}).
					Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id AND manager_schedule.deleted_at IS NULL").
					Where("manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ? AND building.deleted_at IS NULL", *managerID).Select("building.id")).
			Limit(int(limit)).Offset(int(offset)).Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
			Find(tickets).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SupportTicketRepository) GetSupportTicket(ctx *gin.Context, ticket *structs.SupportTicket, ticketID int64, managerID *int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
			Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
			Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
			Where("support_ticket.id = ? AND support_ticket.manager_id IS NOT NULL", ticketID).
			Find(ticket).Error; err != nil {
			return err
		}
	} else {
		if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
			Joins("JOIN contract ON contract.id = support_ticket.contract_id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("support_ticket.id = ? AND (manager_id = ? OR (manager_id IS NULL AND building.id IN (?))) AND support_ticket.deleted_at IS NULL", ticketID, *managerID,
				config.DB.Model(&models.BuildingModel{}).
					Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id AND manager_schedule.deleted_at IS NULL").
					Where("manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ? AND building.deleted_at IS NULL", *managerID).Select("building.id")).
			Find(ticket).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketsByManagerID(ctx *gin.Context, tickets *[]structs.SupportTicket, managerID int64, limit int64, offset int64, startDate string, endDate string) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Distinct().Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
		Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
		Where("support_ticket.manager_id = ? AND support_ticket.created_at::timestamp::date >= ? AND support_ticket.created_at::timestamp::date <= ? AND support_ticket.deleted_at IS NULL", managerID, startDate, endDate).
		Limit(int(limit)).Offset(int(offset)).Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
		Find(tickets).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketsByCustomerID(ctx *gin.Context, tickets *[]structs.SupportTicket, customerID int64) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Distinct().Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
		Joins("LEFT JOIN room_resident_list ON room_resident_list.contract_id = contract.id").
		Joins("JOIN room_resident ON room_resident_list.resident_id = room_resident.id AND room_resident.deleted_at IS NULL").
		Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
		Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
		Where("(contract.householder_id = ? OR room_resident.user_account_id = ?) AND support_ticket.deleted_at IS NULL", customerID, customerID).
		Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
		Find(tickets).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketByCustomerID(ctx *gin.Context, ticket *structs.SupportTicket, customerID int64, ticketID int64) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Distinct().Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
		Joins("LEFT JOIN room_resident_list ON room_resident_list.contract_id = contract.id").
		Joins("JOIN room_resident ON room_resident_list.resident_id = room_resident.id AND room_resident.deleted_at IS NULL").
		Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
		Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
		Where("(contract.householder_id = ? OR room_resident.user_account_id = ?) AND support_ticket.deleted_at IS NULL AND support_ticket.id = ?", customerID, customerID, ticketID).
		Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
		Find(ticket).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketsByCustomerID2(ctx *gin.Context, tickets *[]structs.SupportTicket, limit, offset int64, startDate string, endDate string, customerID int64) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").Distinct().Select("support_ticket.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON support_ticket.contract_id = contract.id AND contract.deleted_at IS NULL").
		Joins("LEFT JOIN room_resident_list ON room_resident_list.contract_id = contract.id").
		Joins("JOIN room_resident ON room_resident_list.resident_id = room_resident.id AND room_resident.deleted_at IS NULL").
		Joins("JOIN room ON contract.room_id = room.id AND room.deleted_at IS NULL").
		Joins("JOIN building ON room.building_id = building.id AND building.deleted_at IS NULL").
		Where("(contract.householder_id = ? OR room_resident.user_account_id = ?) AND support_ticket.created_at::timestamp::date >= ? AND support_ticket.created_at::timestamp::date <= ? AND support_ticket.deleted_at IS NULL", customerID, customerID, startDate, endDate).
		Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
		Limit(int(limit)).Offset(int(offset)).Order("support_ticket.created_at desc, support_ticket.owner_resolve_time desc, support_ticket.manager_resolve_time desc").
		Find(tickets).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupportTicketRepository) GetTicketByRoomIDAndBuildingID(ctx *gin.Context, roomID int64, buildingID int64, startDate string, endDate string, tickets *[]models.SupportTicketModel) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).Preload("Files").Preload("Manager").Preload("Customer").Preload("Owner").
		Joins("JOIN contract ON contract.id = support_ticket.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("support_ticket.created_at::timestamp::date >= ? AND support_ticket.created_at::timestamp::date <= ? AND room.id = ? AND building.id = ? AND support_ticket.deleted_at IS NULL", startDate, endDate, roomID, buildingID).
		Order("support_ticket.created_at desc, owner_resolve_time desc, manager_resolve_time desc").
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

func (r *SupportTicketRepository) Update(ctx *gin.Context, tx *gorm.DB, ticket *models.SupportTicketModel, id int64) error {
	userID := ctx.GetInt64("userID")
	query := tx.Set("userID", userID).Model(&models.SupportTicketModel{})

	if ticket.OwnerID == 0 {
		query = query.Omit("OwnerID")
	}

	if ticket.ManagerID == 0 {
		query = query.Omit("ManagerID")
	}

	if err := query.Where("id = ?", id).Updates(ticket).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupportTicketRepository) GetDeletableTickets(ctx *gin.Context, tickets *[]models.SupportTicketModel, ids []int64, customerID int64) error {
	if err := config.DB.Model(&models.SupportTicketModel{}).
		Where("id IN ? AND status = ? AND manager_id IS NULL AND owner_id IS NULL AND deleted_at IS NULL AND customer_id = ?", ids, constants.Common.SupportTicketStatus.PENDING, customerID).
		Find(tickets).Error; err != nil {
		return err
	}
	return nil
}

func (r *SupportTicketRepository) DeleteTicketFiles(ctx *gin.Context, tx *gorm.DB, ticketID int64, fileIDs []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.SupportTicketFileModel{}).Where("id in ? AND support_ticket_id = ?", fileIDs, ticketID).UpdateColumns(models.SupportTicketFileModel{
		DefaultFileModel: models.DefaultFileModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *SupportTicketRepository) AddFile(ctx *gin.Context, tx *gorm.DB, images *[]models.SupportTicketFileModel) error {
	userID := ctx.GetInt64("userID")
	// for _, image := range *images {
	if err := tx.Set("userID", userID).Model(&models.SupportTicketFileModel{}).Omit("ID").Save(&images).Error; err != nil {
		return err
		// }
	}
	return nil
}

func (r *SupportTicketRepository) Add(ctx *gin.Context, tx *gorm.DB, ticket *models.SupportTicketModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.SupportTicketModel{}).Omit("ID", "ManagerID", "OwnerID").Create(ticket).Error; err != nil {
		return err
	}

	return nil
}
