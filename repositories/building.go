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

type BuildingRepository struct {
}

func NewBuildingRepository() *BuildingRepository {
	return &BuildingRepository{}
}

func (r *BuildingRepository) Get(ctx *gin.Context, building *[]models.BuildingModel) error {
	if err := config.DB.Model(&models.BuildingModel{}).Preload("Images").Preload("Rooms").Order("id asc").Find(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetBuildingBaseOnSchedule(ctx *gin.Context, building *[]models.BuildingModel, userID int64) error {
	if err := config.DB.Model(&models.BuildingModel{}).Preload("Images").Preload("Rooms").
		Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id").
		Where("manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ? AND building.deleted_at IS NULL AND manager_schedule.deleted_at IS NULL", userID).Order("id asc").
		Find(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetById(ctx *gin.Context, building *models.BuildingModel, id int64) error {
	if err := config.DB.Model(&models.BuildingModel{}).Where("id = ?", id).Preload("Rooms").Preload("Rooms.Images").Preload("Images").Preload("Services").Find(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetNewImageNo(ctx *gin.Context, buildingID int64) (int, error) {
	lastestImage := models.BuildingImageModel{}
	if err := config.DB.Model(&models.BuildingImageModel{}).Where("building_id = ?", buildingID).Order("no desc").Unscoped().Find(&lastestImage).Error; err != nil {
		return 0, err
	}
	return lastestImage.No + 1, nil
}

func (r *BuildingRepository) Create(ctx *gin.Context, tx *gorm.DB, building *models.BuildingModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Model(&models.BuildingModel{}).Omit("ID").Create(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) Update(ctx *gin.Context, tx *gorm.DB, building *models.BuildingModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Save(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.BuildingModel{}).Where("id in ?", id).UpdateColumns(models.BuildingModel{
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

	// if err := tx.Set("isQuiet", true).Model(&models.BuildingImageModel{}).Where("building_id in ?", id).UpdateColumns(models.BuildingImageModel{
	// 	DefaultFileModel: models.DefaultFileModel{
	// 		DeletedBy: userID,
	// 		DeletedAt: gorm.DeletedAt{
	// 			Valid: true,
	// 			Time:  now,
	// 		},
	// 	},
	// }).Error; err != nil {
	// 	return err
	// }

	// if err := tx.Set("isQuiet", true).Model(&models.BuildingServiceModel{}).Where("building_id in ?", id).UpdateColumns(models.BuildingServiceModel{
	// 	DefaultModel: models.DefaultModel{
	// 		DeletedBy: userID,
	// 		DeletedAt: gorm.DeletedAt{
	// 			Valid: true,
	// 			Time:  now,
	// 		},
	// 	},
	// }).Error; err != nil {
	// 	return err
	// }
	return nil
}

func (r *BuildingRepository) AddImage(ctx *gin.Context, tx *gorm.DB, image *[]models.BuildingImageModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BuildingImageModel{}).Omit("ID").Create(image).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) DeleteImages(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.BuildingImageModel{}).Where("id in ?", id).UpdateColumns(models.BuildingImageModel{
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

func (r *BuildingRepository) AddServices(ctx *gin.Context, tx *gorm.DB, services *[]models.BuildingServiceModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BuildingServiceModel{}).Omit("ID").Create(services).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) DeleteServices(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.BuildingServiceModel{}).Where("id in ?", id).UpdateColumns(models.BuildingServiceModel{
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

func (r *BuildingRepository) EditService(ctx *gin.Context, service *models.BuildingServiceModel) error {
	userID := ctx.GetInt64("userID")
	if err := config.DB.Set("userID", userID).Save(service).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetServiceByID(ctx *gin.Context, service *models.BuildingServiceModel, id int64) error {
	if err := config.DB.Model(&models.BuildingServiceModel{}).
		Joins("JOIN building ON building.id = building_service.building_id AND building.deleted_at IS NULL").
		Where("building_service.id = ? AND building_service.deleted_at IS NULL", id).Find(service).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetBuildingSchedule(ctx *gin.Context, buildingID int64, schedule *[]models.ManagerScheduleModel) error {
	if err := config.DB.Model(&models.ManagerScheduleModel{}).Preload("Manager", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Building").
		Joins("JOIN building on building.id = manager_schedule.building_id").
		Where("building_id = ? AND building.deleted_at IS NULL AND manager_schedule.deleted_at IS NULL", buildingID).Find(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetBuildingRoom(ctx *gin.Context, buildingID int64, rooms *[]models.RoomModel) error {
	if err := config.DB.Model(&models.RoomModel{}).Where("building_id = ?", buildingID).Find(rooms).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetManagerBuildingSchedule(ctx *gin.Context, buildingID int64, schedule *[]models.ManagerScheduleModel, mangerID int64) error {
	if err := config.DB.Model(&models.ManagerScheduleModel{}).Preload("Manager", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Building").
		Joins("JOIN building on building.id = manager_schedule.building_id").
		Where("building_id = ? AND manager_id = ? AND building.deleted_at IS NULL AND manager_schedule.deleted_at IS NULL", buildingID, mangerID).Find(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetServicesByIDs(ctx *gin.Context, services *[]models.BuildingServiceModel, IDs []int64) error {
	if err := config.DB.Model(&models.BuildingServiceModel{}).
		Joins("JOIN building ON building.id = building_service.building_id AND building.deleted_at IS NULL").
		Where("building_service.id in ? AND building_service.deleted_at IS NULL", IDs).Find(services).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) UpdateServices(ctx *gin.Context, tx *gorm.DB, services *[]models.BuildingServiceModel) error {
	userID := ctx.GetInt64("userID")
	// if err := tx.Set("userID", userID).Save(services).Error; err != nil {
	// 	return err
	// }

	for _, service := range *services {
		if err := tx.Set("userID", userID).Model(&models.BuildingServiceModel{}).Where("id = ?", service.ID).Save(service).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *BuildingRepository) GetBuildingByContractID(ctx *gin.Context, building *models.BuildingModel, ID int64) error {
	if err := config.DB.Model(&models.BuildingModel{}).
		Joins("JOIN room ON room.building_id = building.id AND room.deleted_at IS NULL").
		Joins("JOIN contract ON contract.room_id = room.id AND contract.deleted_at IS NULL").
		Where("contract.id = ? AND building.deleted_at IS NULL", ID).Find(building).Error; err != nil {
		return err
	}

	return nil
}

func (r *BuildingRepository) GetAllBuildingStatistic(ctx *gin.Context, data *structs.AllBuildingStatistic) error {
	if err := config.DB.Raw(`SELECT
		(SELECT COUNT(*) FROM building WHERE deleted_at IS NULL) as total_buildings,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.deleted_at IS NULL) as total_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL) as total_rented_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL) as total_bought_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL) as total_available_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL) as total_maintenanced_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL) as total_unavailable_rooms
	`, constants.Common.RoomStatus.RENTED, constants.Common.RoomStatus.SOLD, constants.Common.RoomStatus.AVAILABLE, constants.Common.RoomStatus.MAINTENANCE, constants.Common.RoomStatus.UNAVAILABLE).Scan(data).Error; err != nil {
		return err
	}

	return nil
}

func (r *BuildingRepository) GetBuildingStatistic(ctx *gin.Context, buildingID int64, data *structs.BuildingStatistic, year int64) error {
	if err := config.DB.Raw(`SELECT
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.deleted_at IS NULL AND building.id = ?) as total_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL AND building.id = ?) as total_rented_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL AND building.id = ?) as total_bought_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL AND building.id = ?) as total_available_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL AND building.id = ?) as total_maintenanced_rooms,
		(SELECT COUNT(*) FROM room JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE room.status = ? AND room.deleted_at IS NULL AND building.id = ?) as total_unavailable_rooms,

		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND building.id = ?) AS total_contract,
		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND building.id = ?) AS total_rent_contract,
		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND building.id = ?) AS total_buy_contract,

		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_active_rent_contract,
		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_active_buy_contract,

		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_expire_rent_contract,

		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_cancel_rent_contract,
		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_cancel_buy_contract,

		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_wait_for_signature_rent_contract,
		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_wait_for_signature_buy_contract,

		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_not_in_effect_rent_contract,
		(SELECT COUNT(*) FROM contract JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE contract.deleted_at IS NULL AND contract.type = ? AND contract.status = ? AND building.id = ?) AS total_not_in_effect_buy_contract,

		(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND building.id = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_bill,
		(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND bill.status = ? AND building.id = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_paid,
		(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND bill.status = ? AND building.id = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_unpaid,
		(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND bill.status = ? AND building.id = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_overdue
	`, buildingID,
		constants.Common.RoomStatus.RENTED,
		buildingID,
		constants.Common.RoomStatus.SOLD,
		buildingID,
		constants.Common.RoomStatus.AVAILABLE,
		buildingID,
		constants.Common.RoomStatus.MAINTENANCE,
		buildingID,
		constants.Common.RoomStatus.UNAVAILABLE,
		buildingID,
		buildingID,
		constants.Common.ContractType.RENT, buildingID, constants.Common.ContractType.BUY, buildingID,
		constants.Common.ContractType.RENT, constants.Common.ContractStatus.ACTIVE, buildingID, constants.Common.ContractType.BUY, constants.Common.ContractStatus.ACTIVE, buildingID,
		constants.Common.ContractType.RENT, constants.Common.ContractStatus.EXPIRED, buildingID,
		constants.Common.ContractType.RENT, constants.Common.ContractStatus.CANCELLED, buildingID, constants.Common.ContractType.BUY, constants.Common.ContractStatus.CANCELLED, buildingID,
		constants.Common.ContractType.RENT, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE, buildingID, constants.Common.ContractType.BUY, constants.Common.ContractStatus.WAITING_FOR_SIGNATURE, buildingID,
		constants.Common.ContractType.RENT, constants.Common.ContractStatus.NOT_IN_EFFECT, buildingID, constants.Common.ContractType.BUY, constants.Common.ContractStatus.NOT_IN_EFFECT, buildingID,
		buildingID,
		constants.Common.BillStatus.PAID, buildingID, constants.Common.BillStatus.UN_PAID, buildingID, constants.Common.BillStatus.OVERDUE, buildingID,
	).Scan(data).Error; err != nil {
		return err
	}

	data.RevenueStatistic = []structs.RevenueStatisticStruct{}

	if err := config.DB.Raw(`
    SELECT
        b.period AS period,
        SUM(b.amount) AS total_expected_revenue,
        SUM(CASE WHEN b.status IN ? THEN b.amount ELSE 0 END) AS total_actual_revenue,
        SUM(CASE WHEN b.status NOT IN ? THEN b.amount ELSE 0 END) AS total_remaining_revenue
    FROM bill AS b JOIN contract ON contract.id = b.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL
    WHERE b.deleted_at IS NULL AND building.id = ? AND EXTRACT(YEAR FROM b.period) = ?
    GROUP BY period
    ORDER BY period
    `,
		[]int{constants.Common.BillStatus.PAID},
		[]int{constants.Common.BillStatus.PAID},
		buildingID, year,
	).Scan(&data.RevenueStatistic).Error; err != nil {
		return err
	}

	return nil
}
