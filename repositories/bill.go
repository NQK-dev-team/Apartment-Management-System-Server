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

type BillRepository struct {
}

func NewBillRepository() *BillRepository {
	return &BillRepository{}
}

func (r *BillRepository) GetById(ctx *gin.Context, bill *structs.Bill, id int64) error {
	if err := config.DB.Model(&models.BillModel{}).Preload("Contract").Preload("Contract.Householder", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").Select("bill.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("bill.id = ? AND bill.deleted_at IS NULL", id).Find(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetById2(ctx *gin.Context, bill *models.BillModel, id int64) error {
	if err := config.DB.Model(&models.BillModel{}).Preload("Contract").Preload("Contract.Householder", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("bill.id = ? AND bill.deleted_at IS NULL", id).Find(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetBillList(ctx *gin.Context, bills *[]structs.Bill, startMonth, endMonth string, limit, offset int64) error {
	if err := config.DB.Model(&models.BillModel{}).Preload("Contract").Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").Select("bill.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("bill.deleted_at IS NULL AND period BETWEEN ? AND ?", startMonth, endMonth).Order("payment_time DESC").
		Limit(int(limit)).Offset(int(offset)).
		Find(bills).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetBillListForManager(ctx *gin.Context, bills *[]structs.Bill, startMonth, endMonth string, limit, offset int64, managerID int64) error {
	if err := config.DB.Model(&models.BillModel{}).Preload("Contract").Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").Select("bill.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id AND manager_schedule.deleted_at IS NULL").
		Where("bill.deleted_at IS NULL AND period BETWEEN ? AND ? AND manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ?", startMonth, endMonth, managerID).
		Order("payment_time DESC").Limit(int(limit)).Offset(int(offset)).
		Find(bills).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetBillListForCustomer(ctx *gin.Context, bills *[]structs.Bill, startMonth, endMonth string, limit, offset int64, customerID int64) error {
	contractQuery := config.DB.Model(&models.ContractModel{}).Select("contract.id").Distinct("contract.id").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Joins("LEFT JOIN room_resident_list ON room_resident_list.contract_id = contract.id").
		Joins("JOIN room_resident ON room_resident.id = room_resident_list.resident_id AND room_resident.deleted_at IS NULL").
		Where("contract.deleted_at IS NULL AND (contract.householder_id = ? OR room_resident.user_account_id = ?)", customerID, customerID)

	if err := config.DB.Model(&models.BillModel{}).Preload("Contract").Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").Select("bill.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("bill.deleted_at IS NULL AND period BETWEEN ? AND ? AND contract.id IN (?)", startMonth, endMonth, contractQuery).
		Order("payment_time DESC").Limit(int(limit)).Offset(int(offset)).
		Find(bills).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetByContractId(ctx *gin.Context, bills *[]models.BillModel, contractID int64) error {
	if err := config.DB.Model(&models.BillModel{}).
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("contract_id = ? AND bill.deleted_at IS NULL", contractID).Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").Find(bills).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.BillModel{}).Where("id IN ?", id).UpdateColumns(models.BillModel{
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

	if err := tx.Set("isQuiet", true).Model(&models.BillPaymentModel{}).Where("bill_id IN ?", id).UpdateColumns(models.BillPaymentModel{
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

func (r *BillRepository) GetDeletableBills(ctx *gin.Context, bills *[]models.BillModel, IDs []int64, managerID *int64) error {
	if managerID == nil {
		if err := config.DB.Model(&models.BillModel{}).
			Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("bill.id in ? and bill.status in ? AND bill.deleted_at IS NULL", IDs, []int{constants.Common.BillStatus.UN_PAID, constants.Common.BillStatus.OVERDUE}).Find(bills).Error; err != nil {
			return err
		}
	} else {
		query1 := config.DB.Model(&models.BillModel{}).
			Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Where("bill.id in ? and bill.status in ? and contract.creator_id = ? AND bill.deleted_at IS NULL", IDs, []int{constants.Common.BillStatus.UN_PAID, constants.Common.BillStatus.OVERDUE}, *managerID)
		query2 := config.DB.Model(&models.BillModel{}).
			Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
			Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
			Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
			Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id AND manager_schedule.deleted_at IS NULL").
			Where("bill.id in ? and bill.status in ? and contract.creator_id != ? AND bill.deleted_at IS NULL AND manager_schedule.start_date <= now() AND COALESCE(manager_schedule.end_date,now()) >= now() AND manager_schedule.manager_id = ?", IDs, []int{constants.Common.BillStatus.UN_PAID, constants.Common.BillStatus.OVERDUE}, *managerID, *managerID)

		if err := config.DB.Model(&models.BillModel{}).Table("((?) UNION ALL (?)) as all_bills", query1, query2).
			Find(bills).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *BillRepository) GetBillBuildingID(ctx *gin.Context, billID int64, buildingID *int64) error {
	if err := config.DB.Model(&models.BillModel{}).
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("bill.id = ? AND bill.deleted_at IS NULL", billID).Select("building.id").Scan(buildingID).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetBillByIDForCustomer(ctx *gin.Context, bill *structs.Bill, customerID, billID int64) error {
	contractQuery := config.DB.Model(&models.ContractModel{}).Select("contract.id").Distinct("contract.id").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Joins("LEFT JOIN room_resident_list ON room_resident_list.contract_id = contract.id").
		Joins("JOIN room_resident ON room_resident.id = room_resident_list.resident_id AND room_resident.deleted_at IS NULL").
		Where("contract.deleted_at IS NULL AND (contract.householder_id = ? OR room_resident.user_account_id = ?)", customerID, customerID)

	if err := config.DB.Model(&models.BillModel{}).Preload("Contract").Preload("Payer", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("BillPayments").Select("bill.*, building.name AS building_name, room.no AS room_no, room.floor AS room_floor").
		Joins("JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL").
		Joins("JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL").
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("bill.deleted_at IS NULL AND bill.id = ? AND contract.id IN (?)", billID, contractQuery).
		Order("payment_time DESC").
		Find(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) DeletePayment(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.BillPaymentModel{}).Where("id in ?", id).UpdateColumns(models.BillPaymentModel{
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

func (r *BillRepository) AddNewPayment(ctx *gin.Context, tx *gorm.DB, payment *[]models.BillPaymentModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BillPaymentModel{}).Omit("ID").Save(payment).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) UpdatePayment(ctx *gin.Context, tx *gorm.DB, payment *[]models.BillPaymentModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BillPaymentModel{}).Omit("CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy").Save(payment).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) UpdateBill(ctx *gin.Context, tx *gorm.DB, bill *models.BillModel, id int64) error {
	userID := ctx.GetInt64("userID")
	if bill.PayerID.Int64 == 0 || !bill.PayerID.Valid {
		if err := tx.Set("userID", userID).Model(&models.BillModel{}).Omit("PayerID").Where("id = ?", id).Save(bill).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Set("userID", userID).Model(&models.BillModel{}).Where("id = ?", id).Save(bill).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *BillRepository) UpdateBill2(tx *gorm.DB, bill *models.BillModel, id int64) error {
	if bill.PayerID.Int64 == 0 || !bill.PayerID.Valid {
		if err := tx.Model(&models.BillModel{}).Omit("PayerID").Where("id = ?", id).Save(bill).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(&models.BillModel{}).Where("id = ?", id).Save(bill).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *BillRepository) CancelBillPayment(tx *gorm.DB, billID int64) error {
	if err := tx.Model(&models.BillModel{}).Select("payer_id", "payment_time", "status", "order_id", "request_id").Where("id = ?", billID).Updates(map[string]interface{}{
		"payer_id":     nil,
		"payment_time": nil,
		"status":       constants.Common.BillStatus.UN_PAID,
		"order_id":     nil,
		"request_id":   nil,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) CreateBill(ctx *gin.Context, tx *gorm.DB, bill *models.BillModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BillModel{}).Omit("ID").Create(bill).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) UpdateBillOfCancelContract(ctx *gin.Context, tx *gorm.DB, contractID int64) error {
	cancelTime := time.Now().Format("2006-01-02")

	if err := tx.Exec("UPDATE bill SET status = ? WHERE payer_id IS NULL AND payment_time IS NULL AND deleted_at IS NULL AND status = ? AND contract_id = ? AND period <= ?", constants.Common.BillStatus.OVERDUE, constants.Common.BillStatus.UN_PAID, contractID, cancelTime).Error; err != nil {
		return err
	}

	if err := tx.Exec("UPDATE bill SET status = ? WHERE payer_id IS NULL AND payment_time IS NULL AND deleted_at IS NULL AND status = ? AND contract_id = ? AND period > ?", constants.Common.BillStatus.CANCELLED, constants.Common.BillStatus.UN_PAID, contractID, cancelTime).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) UpdateBillStatus(tx *gorm.DB) error {
	if err := tx.Exec("UPDATE bill SET status = ? WHERE payer_id IS NULL AND payment_time IS NULL AND deleted_at IS NULL AND date_trunc('month', period) < date_trunc('month', NOW()) AND status = ? AND contract_id IN (SELECT contract.id FROM contract where contract.deleted_at IS NULL AND contract.status = ?)", constants.Common.BillStatus.OVERDUE, constants.Common.BillStatus.UN_PAID, constants.Common.ContractStatus.ACTIVE).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
	UPDATE bill SET status =
		CASE
			WHEN bill.period <= contract.end_date THEN ?
			WHEN bill.period > contract.end_date THEN ?
			ELSE bill.status
		END
	FROM contract
	WHERE bill.contract_id = contract.id AND bill.payer_id IS NULL AND bill.payment_time IS NULL AND bill.deleted_at IS NULL AND contract.deleted_at IS NULL AND bill.status = ? AND contract.status = ?
		`, constants.Common.BillStatus.OVERDUE, constants.Common.BillStatus.CANCELLED, constants.Common.ContractStatus.ACTIVE, constants.Common.ContractStatus.EXPIRED).Error; err != nil {
		return err
	}

	return nil
}

func (r *BillRepository) AddNewPayment2(ctx *gin.Context, tx *gorm.DB, payment *models.BillPaymentModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BillPaymentModel{}).Omit("ID").Save(payment).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetMomoBills(bills *[]models.BillModel) error {
	if err := config.DB.Model(&models.BillModel{}).Where("status = ?", constants.Common.BillStatus.PROCESSING).Find(bills).Error; err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetBillStatistics(ctx *gin.Context, data *structs.BillStatistic, year int64) error {
	if err := config.DB.Raw(`SELECT
	(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_bill,
	(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND bill.status = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_paid,
	(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND bill.status = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_unpaid,
	(SELECT COUNT(*) FROM bill JOIN contract ON contract.id = bill.contract_id AND contract.deleted_at IS NULL JOIN room ON room.id = contract.room_id AND room.deleted_at IS NULL JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL WHERE bill.deleted_at IS NULL AND bill.status = ? AND date_trunc('month', bill.period) = date_trunc('month', NOW())) as total_overdue
	`, constants.Common.BillStatus.PAID, constants.Common.BillStatus.UN_PAID, constants.Common.BillStatus.OVERDUE,
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
    WHERE b.deleted_at IS NULL AND EXTRACT(YEAR FROM b.period) = ?
    GROUP BY period
    ORDER BY period
    `,
		[]int{constants.Common.BillStatus.PAID},
		[]int{constants.Common.BillStatus.PAID},
		year,
	).Scan(&data.RevenueStatistic).Error; err != nil {
		return err
	}
	return nil
}
