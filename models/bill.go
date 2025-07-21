package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BillModel struct {
	DefaultModel
	Title        string             `json:"title" gorm:"column:title;not null;type:varchar(255);"`
	Period       time.Time          `json:"period" gorm:"column:period;not null;type:date;"`
	Status       int                `json:"status" gorm:"column:status;not null;type:int;"` // 1: Unpaid, 2: Paid, 3: Overdue, 4: Processing, 5: Cancelled
	Note         sql.NullString     `json:"note" gorm:"column:note;type:varchar(255);"`
	PaymentTime  sql.NullTime       `json:"paymentTime" gorm:"column:payment_time;type:timestamp;"`
	Amount       float64            `json:"amount" gorm:"column:amount;not null;type:numeric;"`
	PayerID      sql.NullInt64      `json:"payerID" gorm:"column:payer_id;"`
	Payer        UserModel          `json:"payer" gorm:"foreignKey:payer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContractID   int64              `json:"contractID" gorm:"column:contract_id;not null;"`
	BillPayments []BillPaymentModel `json:"billPayments" gorm:"foreignKey:bill_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Contract     ContractModel      `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (u *BillModel) TableName() string {
	return "bill"
}

func (u *BillModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *BillModel) BeforeUpdate(tx *gorm.DB) error {
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
