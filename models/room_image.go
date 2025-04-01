package models

import (
	"time"

	"gorm.io/gorm"
)

type RoomImageModel struct {
	DefaultFileModel
	RoomID int64 `json:"roomID" gorm:"column:room_id;not null;"`
	// BuildingID int64 `json:"buildingID" gorm:"column:building_id;not null;"`
	// Room       RoomModel `json:"room" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *RoomImageModel) TableName() string {
	return "room_image"
}

func (u *RoomImageModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
	}
	u.CreatedAt = time.Now()

	return nil
}
