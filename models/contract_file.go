package models

import (
	"time"

	"gorm.io/gorm"
)

type ContractFileModel struct {
	DefaultFileModel
	ContractID int64 `json:"contractID" gorm:"column:contract_id;not null;"`
	// Contract   ContractModel `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *ContractFileModel) TableName() string {
	return "contract_file"
}

func (u *ContractFileModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
	}
	u.CreatedAt = time.Now()

	return nil
}
