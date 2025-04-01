package models

import (
	"time"

	"gorm.io/gorm"
)

type BuildingImageModel struct {
	DefaultFileModel
	BuildingID int64 `json:"buildingID" gorm:"column:building_id;not null;"`
	// Building   BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *BuildingImageModel) TableName() string {
	return "building_image"
}

func (u *BuildingImageModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
	}
	u.CreatedAt = time.Now()

	return nil
}
