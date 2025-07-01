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
	if err := config.DB.Model(&models.ContractModel{}).Where("id = ? AND contract.deleted_at IS NULL", id).Preload("Files").Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByIDs(ctx *gin.Context, contract *[]models.ContractModel, id []int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Where("id in ? AND contract.deleted_at IS NULL", id).Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByRoomID(ctx *gin.Context, contract *[]models.ContractModel, roomIDs []int64) error {
	if err := config.DB.Model(&models.ContractModel{}).
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("contract.room_id in ? AND contract.deleted_at IS NULL", roomIDs).Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractsByManagerID(ctx *gin.Context, contracts *[]models.ContractModel, managerID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").
		Where("creator_id = ? AND contract.deleted_at IS NULL", managerID).Order("start_date DESC, end_date DESC, sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByID(ctx *gin.Context, contract *structs.Contract, id int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Preload("Files").Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("contract.id = ? AND contract.deleted_at IS NULL", id).
		Find(contract).Error; err != nil {
		return err
	}

	// if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contract).ID).Scan(&(*contract).RoomNo).Error; err != nil {
	// 	return err
	// }

	// if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contract).ID).Scan(&(*contract).RoomFloor).Error; err != nil {
	// 	return err
	// }

	// if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contract).ID).Scan(&(*contract).BuildingName).Error; err != nil {
	// 	return err
	// }

	// if err := config.DB.Raw("SELECT building.address AS building_address FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contract).ID).Scan(&(*contract).BuildingAddress).Error; err != nil {
	// 	return err
	// }

	if err := config.DB.Model(&models.RoomResidentModel{}).Preload("UserAccount").
		Joins("JOIN room_resident_list ON room_resident_list.resident_id = room_resident.ID").
		Where("room_resident_list.contract_id = ? AND room_resident.deleted_at IS NULL", (*contract).ID).
		Find(&(*contract).Residents).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	(*contract).Residents = []models.RoomResidentModel{}
		// } else {
		// 	return err
		// }
	}

	return nil
}

func (r *ContractRepository) GetContracts(ctx *gin.Context, contracts *[]structs.Contract, limit int64, offset int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("contract.deleted_at IS NULL").
		Limit(int(limit)).Offset(int(offset)).Order("contract.start_date DESC, contract.end_date DESC, contract.sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	// for i := range *contracts {
	// 	if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (r *ContractRepository) GetContractsByManagerID2(ctx *gin.Context, contracts *[]structs.Contract, managerID int64, limit int64, offset int64) error {
	query1 := config.DB.Model(&models.ContractModel{}).Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("creator_id = ? AND contract.deleted_at IS NULL", managerID)
	query2 := config.DB.Model(&models.ContractModel{}).Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id AND manager_schedule.deleted_at IS NULL").
		Where("creator_id != ? AND contract.deleted_at IS NULL AND manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ?", managerID, managerID)

	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").
		Table("((?) UNION ALL (?)) as all_contracts", query1, query2).
		Select("all_contracts.*").
		Limit(int(limit)).Offset(int(offset)).Order("all_contracts.start_date DESC, all_contracts.end_date DESC, all_contracts.sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	// for i := range *contracts {
	// 	if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (r *ContractRepository) GetContractsByCustomerID(ctx *gin.Context, contracts *[]structs.Contract, customerID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").
		Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("contract.householder_id = ? AND contract.deleted_at IS NULL", customerID).Order("contract.start_date DESC, contract.end_date DESC, contract.sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	// for i := range *contracts {
	// 	if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (r *ContractRepository) GetContractsByCustomerID2(ctx *gin.Context, contracts *[]structs.Contract, customerID int64, limit int64, offset int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").
		Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("contract.householder_id = ? AND contract.deleted_at IS NULL", customerID).Limit(int(limit)).Offset(int(offset)).Order("contract.start_date DESC, contract.end_date DESC, contract.sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	// for i := range *contracts {
	// 	if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (r *ContractRepository) GetContractByRoomIDAndBuildingID(ctx *gin.Context, contracts *[]structs.Contract, roomID int64, buildingID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").
		Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("room_id = ? AND building_id = ? AND contract.deleted_at IS NULL", roomID, buildingID).Order("contract.start_date DESC, contract.end_date DESC, contract.sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	// for i := range *contracts {
	// 	if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (r *ContractRepository) GetContractByRoomIDAndBuildingIDAndManagerID(ctx *gin.Context, contracts *[]structs.Contract, roomID int64, buildingID int64, managerID int64) error {
	if err := config.DB.Model(&models.ContractModel{}).Preload("Creator").Preload("Householder").
		Select("contract.*, room.no AS room_no, room.floor AS room_floor, building.name AS building_name, building.address AS building_address").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("room_id = ? AND building_id = ? AND creator_id = ? AND contract.deleted_at IS NULL", roomID, buildingID, managerID).Order("contract.start_date DESC, contract.end_date DESC, contract.sign_date DESC").
		Find(contracts).Error; err != nil {
		return err
	}

	// for i := range *contracts {
	// 	if err := config.DB.Raw("SELECT room.no AS room_no FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomNo).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT room.floor AS room_floor FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].RoomFloor).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := config.DB.Raw("SELECT building.name AS building_name FROM building JOIN room ON building.id = room.building_id JOIN contract ON contract.room_id = room.id WHERE contract.id = ? AND room.deleted_at IS NULL AND building.deleted_at IS NULL AND contract.deleted_at IS NULL", (*contracts)[i].ID).Scan(&(*contracts)[i].BuildingName).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (r *ContractRepository) GetDeletableContracts(ctx *gin.Context, contracts *[]models.ContractModel, IDs []int64, managerID *int64, roomID int64, buildingID int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("contract.id in ? and contract.status in ? and room_id = ? and building_id = ? AND contract.deleted_at IS NULL", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, roomID, buildingID).Find(contracts).Error; err != nil {
			return err
		}
	} else {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Where("contract.id in ? and contract.status in ? and creator_id = ? and room_id = ? and building_id = ? AND contract.deleted_at IS NULL", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, *managerID, roomID, buildingID).Find(contracts).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *ContractRepository) GetDeletableContracts2(ctx *gin.Context, contracts *[]models.ContractModel, IDs []int64, managerID *int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.ContractModel{}).
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("contract.id in ? and contract.status in ? AND contract.deleted_at IS NULL", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}).Find(contracts).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
	} else {
		query1 := config.DB.Model(&models.ContractModel{}).
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("contract.id in ? and contract.status in ? and contract.creator_id = ? AND contract.deleted_at IS NULL", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, *managerID)
		query2 := config.DB.Model(&models.ContractModel{}).
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id AND manager_schedule.deleted_at IS NULL").
			Where("contract.id in ? and contract.status in ? and contract.creator_id != ? AND contract.deleted_at IS NULL AND manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ?", IDs, []int{constants.Common.ContractStatus.CANCELLED, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE}, *managerID, *managerID)

		if err := config.DB.Model(&models.ContractModel{}).Table("((?) UNION ALL (?)) as all_contracts", query1, query2).
			Find(contracts).Error; err != nil {
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

func (r *ContractRepository) GetNewFileNo(ctx *gin.Context, contractID int64) (int, error) {
	latestFile := models.ContractFileModel{}
	if err := config.DB.Model(&models.ContractFileModel{}).Where("contract_id = ?", contractID).Order("no desc").Unscoped().Find(&latestFile).Error; err != nil {
		return 0, err
	}
	return latestFile.No + 1, nil
}

func (r *ContractRepository) AddFile(ctx *gin.Context, tx *gorm.DB, file *[]models.ContractFileModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.ContractFileModel{}).Omit("ID").Create(file).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) DeleteResident(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.RoomResidentModel{}).Where("id in ?", id).UpdateColumns(models.RoomResidentModel{
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
	return nil
}

func (r *ContractRepository) UpdateContract(ctx *gin.Context, tx *gorm.DB, contract *models.ContractModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.ContractModel{}).Where("id = ?", contract.ID).Save(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) AddNewRoomResident(ctx *gin.Context, tx *gorm.DB, resident *models.RoomResidentModel, contractID int64) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.RoomResidentModel{}).Omit("ID").Create(resident).Error; err != nil {
		return err
	}

	if err := tx.Set("userID", userID).Model(&models.RoomResidentListModel{}).Create(&models.RoomResidentListModel{
		ResidentID: resident.ID,
		ContractID: contractID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) UpdateRoomResident(ctx *gin.Context, tx *gorm.DB, resident *models.RoomResidentModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.RoomResidentModel{}).Where("id = ?", resident.ID).Save(resident).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) UpdateContractStatus(tx *gorm.DB) error {
	if err := tx.Model(&models.ContractModel{}).
		Where("sign_date IS NOT NULL AND sign_date <= now() AND end_date >= now() AND status NOT IN ?", []int{constants.Common.ContractStatus.ACTIVE, constants.Common.ContractStatus.EXPIRED, constants.Common.ContractStatus.CANCELLED}).
		Update("status", constants.Common.ContractStatus.ACTIVE).Error; err != nil {
		return err
	}

	if err := tx.Model(&models.ContractModel{}).
		Where("sign_date IS NOT NULL AND sign_date <= now() AND end_date < now() AND status NOT IN ?", []int{constants.Common.ContractStatus.EXPIRED}).
		Update("status", constants.Common.ContractStatus.EXPIRED).Error; err != nil {
		return err
	}

	if err := tx.Model(&models.ContractModel{}).
		Where("sign_date IS NULL AND sign_date <= now() AND status NOT IN ?", []int{constants.Common.ContractStatus.CANCELLED}).
		Update("status", constants.Common.ContractStatus.CANCELLED).Error; err != nil {
		return err
	}

	return nil
}
