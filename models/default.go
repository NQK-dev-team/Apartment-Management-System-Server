package models

import (
	"time"

	"gorm.io/gorm"
)

type DefaultModel struct {
	ID        int64          `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;type:timestamp;not null;default:now();"`
	CreatedBy int64          `json:"createdBy" gorm:"column:created_by;type:bigint;"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at;type:timestamp;not null;default:now();"`
	UpdatedBy int64          `json:"updatedBy" gorm:"column:updated_by;type:bigint;"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp;"`
	DeletedBy interface{}    `json:"-" gorm:"column:deleted_by;type:bigint;"`
}
