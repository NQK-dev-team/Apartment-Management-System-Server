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

type ContractRepository struct {
}

func NewContractRepository() *ContractRepository {
	return &ContractRepository{}
}

func (r *ContractRepository) GetById(ctx *gin.Context, contract *models.ContractModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("Files").Preload("Bills").Preload("SupportTickets").First(contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByIDs(ctx *gin.Context, contract *[]models.ContractModel, id []int64) error {
	if err := config.DB.Where("id in ?", id).Preload("Bills").Preload("SupportTickets").Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByRoomID(ctx *gin.Context, contract *[]models.ContractModel, roomIDs []int64) error {
	if err := config.DB.Where("room_id in ?", roomIDs).Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractsByManagerID(ctx *gin.Context, contracts *[]models.ContractModel, managerID int64) error {
	if err := config.DB.Preload("Creator").Preload("Householder").Preload("Files").
		Where("creator_id = ?", managerID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractsByCustomerID(ctx *gin.Context, contracts *[]structs.Contract, customerID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Where("householder_id = ?", customerID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	for i := range *contracts {
		if err := config.DB.Raw("SELECT room.no AS room_no FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
			return err
		}

		if err := config.DB.Raw("SELECT building.name AS building_name FROM building INNER JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ?", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *ContractRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.ContractModel{}).Where("id in ?", id).UpdateColumns(models.ContractModel{
		DefaultModel: models.DefaultModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}

	if err := tx.Set("isQuiet", true).Model(&models.ContractFileModel{}).Where("contract_id in ?", id).UpdateColumns(models.ContractFileModel{
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
