package models

import (
	"time"

	"gorm.io/gorm"
)

type DefaultModel struct {
	ID        int64          `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	CreatedBy int64          `json:"createdBy" gorm:"column:created_by;type:bigint;"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at;type:timestamp with time zone;not null;default:now();"`
	UpdatedBy int64          `json:"updatedBy" gorm:"column:updated_by;type:bigint;"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp with time zone;"`
	DeletedBy interface{}    `json:"-" gorm:"column:deleted_by;type:bigint;"`
}

func (u *DefaultModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *DefaultModel) BeforeUpdate(tx *gorm.DB) error {
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
