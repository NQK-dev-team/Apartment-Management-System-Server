package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BillModel struct {
	DefaultModel
	Period        time.Time           `json:"period" gorm:"column:period;not null;type:date;"`
	Status        int                 `json:"status" gorm:"column:status;not null;type:int;"` // 1: Unpaid, 2: Paid, 3: Overdue, 4: Processing
	Note          string              `json:"note" gorm:"column:note;type:varchar(255);"`
	PaymentTime   sql.NullTime        `json:"paymentTime" gorm:"column:payment_time;type:timestamp with time zone;default:now();"`
	Amount        float64             `json:"amount" gorm:"column:amount;not null;type:numeric;"`
	PayerID       int64               `json:"payerID" gorm:"column:payer_id;"`
	Payer         UserModel           `json:"payer" gorm:"foreignKey:payer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContractID    int64               `json:"contractID" gorm:"column:contract_id;not null;"`
	ExtraPayments []ExtraPaymentModel `json:"extraPayments" gorm:"foreignKey:bill_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Contract      ContractModel       `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	// ExtraPayments []ExtraPaymentModel `json:"extraPayments" gorm:"foreignKey:bill_id,contract_id;references:id,contract_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *BillModel) TableName() string {
	return "bill"
}

func (u *BillModel) BeforeDelete(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")

	return tx.Transaction(func(tx1 *gorm.DB) error {
		if err := tx1.Set("userID", userID).Model(&ExtraPaymentModel{}).Where("bill_id = ?", u.ID).Delete(&ExtraPaymentModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}
