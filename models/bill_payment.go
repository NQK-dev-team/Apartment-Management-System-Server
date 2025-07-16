package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BillPaymentModel struct {
	DefaultModel
	Name   string         `json:"name" gorm:"column:name;not null;type:varchar(255);"`
	Amount float64        `json:"amount" gorm:"column:amount;not null;type:numeric;"`
	Note   sql.NullString `json:"note" gorm:"column:note;type:varchar(255);"`
	// ContractID int64     `json:"contractID" gorm:"column:contract_id;not null;"`
	BillID int64 `json:"billID" gorm:"column:bill_id;not null;"`
	// Bill       BillModel `json:"bill" gorm:"foreignKey:bill_id,contract_id;references:id,contract_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *BillPaymentModel) TableName() string {
	return "bill_payment"
}

func (u *BillPaymentModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *BillPaymentModel) BeforeUpdate(tx *gorm.DB) error {
	// if tx.Statement.Changed("UpdatedAt", "UpdatedBy") {
	// 	return errors.New(config.GetMessageCode("CONCURRENCY_ERROR"))
	// }

	isQuiet, _ := tx.Get("isQuiet")
	if isQuiet != nil && isQuiet.(bool) {
		return nil
	}

	userID, _ := tx.Get("userID")
	if userID != nil {
		tx.Statement.SetColumn("updated_by", userID.(int64))
	}
	tx.Statement.SetColumn("updated_at", time.Now())

	return nil
}
