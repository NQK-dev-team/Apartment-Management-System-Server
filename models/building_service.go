package models

import (
	"time"

	"gorm.io/gorm"
)

type BuildingServiceModel struct {
	DefaultModel
	BuildingID int64 `json:"buildingID" gorm:"column:building_id;not null;"`
	// Building   BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name  string  `json:"name" gorm:"column:name;type:varchar(255);not null;"`
	Price float64 `json:"price" gorm:"column:price;type:numeric;not null;"`
}

func (u *BuildingServiceModel) TableName() string {
	return "building_service"
}

func (u *BuildingServiceModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *BuildingServiceModel) BeforeUpdate(tx *gorm.DB) error {
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
