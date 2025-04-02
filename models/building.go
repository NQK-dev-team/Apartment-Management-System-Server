package models

import (
	"time"

	"gorm.io/gorm"
)

type BuildingModel struct {
	DefaultModel
	Name       string                 `json:"name" gorm:"column:name;type:varchar(255);not null;"`
	Address    string                 `json:"address" gorm:"column:address;type:varchar(255);not null;"`
	TotalFloor int                    `json:"totalFloor" gorm:"column:total_floor;type:int;not null;"`
	TotalRoom  int                    `json:"totalRoom" gorm:"column:total_room;type:int;not null;"`
	Images     []BuildingImageModel   `json:"images" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Services   []BuildingServiceModel `json:"services" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rooms      []RoomModel            `json:"rooms" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *BuildingModel) TableName() string {
	return "building"
}

func (u *BuildingModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *BuildingModel) BeforeUpdate(tx *gorm.DB) error {
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
