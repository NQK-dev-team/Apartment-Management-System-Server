package models

import (
	"api/config"
	"errors"
	"time"

	"gorm.io/gorm"
)

type DefaultModel struct {
	ID        int64          `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	CreatedBy string         `json:"createdBy" gorm:"column:created_by;type:varchar(16);not null;"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at;type:timestamp with time zone;not null;default:now();"`
	UpdatedBy string         `json:"updatedBy" gorm:"column:updated_by;type:varchar(16);not null;"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;type:timestamp with time zone;"`
	DeletedBy string         `json:"deletedBy" gorm:"column:deleted_by;type:varchar(16);"`
}

func (u *DefaultModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID == nil {
		userID = "SYSTEM"
	}
	u.CreatedBy = userID.(string)
	u.UpdatedBy = userID.(string)
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *DefaultModel) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("updated_at", "updated_by") {
		return errors.New(config.GetMessageCode("CONCURRENCY_ERROR"))
	}

	userID, _ := tx.Get("userID")
	if userID == nil {
		userID = "SYSTEM"
	}
	tx.Statement.SetColumn("updated_by", userID.(string))
	tx.Statement.SetColumn("updated_at", time.Now())

	return nil
}
