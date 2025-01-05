package models

import (
	"time"

	"gorm.io/gorm"
)

type DefaultModel struct {
	ID        int64          `json:"id" gorm:"primaryKey; colunn:id; autoIncrement; not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null"`
	CreatedBy string         `json:"createdBy" gorm:"column:created_by;type:varchar(16);not null"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at;type:timestamp with time zone;not null"`
	UpdatedBy string         `json:"updatedBy" gorm:"column:updated_by;type:varchar(16);not null"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;type:timestamp with time zone;"`
	DeletedBy string         `json:"deletedBy" gorm:"column:deleted_by;type:varchar(16);"`
}
