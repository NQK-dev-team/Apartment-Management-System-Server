package models

import (
	"time"

	"gorm.io/gorm"
)

type DefaultFileModel struct {
	ID        int64          `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	No        int            `json:"no" gorm:"column:no;type:int;"`
	Title     string         `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Path      string         `json:"path" gorm:"column:path;type:varchar(255);not null;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	CreatedBy int64          `json:"createdBy" gorm:"column:created_by;type:bigint;"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp with time zone;"`
	DeletedBy interface{}    `json:"-" gorm:"column:deleted_by;type:bigint;"`
}

func (u *DefaultFileModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
	}
	u.CreatedAt = time.Now()

	return nil
}
