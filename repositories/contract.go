package repositories

import (
	"api/config"
	"api/constants"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByRoomID(ctx *gin.Context, contract *[]models.ContractModel, roomIDs []int64) error {
	if err := config.DB.Where("room_id in ?", roomIDs).Find(contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractsByManagerID(ctx *gin.Context, contracts *[]models.ContractModel, managerID int64) error {
	if err := config.DB.Preload("Creator").Preload("Householder").Preload("Files").
		Where("creator_id = ?", managerID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *ContractRepository) GetContracts(ctx *gin.Context, contracts *[]structs.Contract, limit int64, offset int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Limit(int(limit)).Offset(int(offset)).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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

func (r *ContractRepository) GetContractsByManagerID2(ctx *gin.Context, contracts *[]structs.Contract, managerID int64, limit int64, offset int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Where("creator_id = ?", managerID).Limit(int(limit)).Offset(int(offset)).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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

func (r *ContractRepository) GetContractsByCustomerID(ctx *gin.Context, contracts *[]structs.Contract, customerID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Where("householder_id = ?", customerID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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

func (r *ContractRepository) GetContractsByCustomerID2(ctx *gin.Context, contracts *[]structs.Contract, customerID int64, limit int64, offset int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Where("householder_id = ?", customerID).Limit(int(limit)).Offset(int(offset)).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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

func (r *ContractRepository) GetContractByRoomIDAndBuildingID(ctx *gin.Context, contracts *[]structs.Contract, roomID int64, buildingID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Joins("INNER JOIN room ON room.id = contract.room_id").
		Where("room_id = ? AND building_id = ?", roomID, buildingID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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

func (r *ContractRepository) GetContractByRoomIDAndBuildingIDAndManagerID(ctx *gin.Context, contracts *[]structs.Contract, roomID int64, buildingID int64, managerID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").
		Joins("INNER JOIN room ON room.id = contract.room_id").
		Where("room_id = ? AND building_id = ? AND creator_id = ?", roomID, buildingID, managerID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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

func (r *ContractRepository) GetDeletableContracts(ctx *gin.Context, contracts *[]models.ContractModel, IDs []int64, managerID *int64, roomID int64, buildingID int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("INNER JOIN room ON room.id = contract.room_id").
			Where("contract.id in ? and contract.status in ? and room_id = ? and building_id = ?", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, roomID, buildingID).Find(contracts).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
	} else {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("INNER JOIN room ON room.id = contract.room_id").
			Where("contract.id in ? and contract.status in ? and creator_id = ? and room_id = ? and building_id = ?", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, *managerID, roomID, buildingID).Find(contracts).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
	}

	return nil
}

func (r *ContractRepository) GetDeletableContracts2(ctx *gin.Context, contracts *[]models.ContractModel, IDs []int64, managerID *int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("INNER JOIN room ON room.id = contract.room_id").
			Where("contract.id in ? and contract.status in ?", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}).Find(contracts).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
	} else {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("INNER JOIN room ON room.id = contract.room_id").
			Where("contract.id in ? and contract.status in ? and creator_id = ?", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, *managerID).Find(contracts).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
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
