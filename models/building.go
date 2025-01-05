package models

import (
	"api/config"
	"errors"
	"time"

	"gorm.io/gorm"
)

type BuildingModel struct {
	DefaultModel
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Address    string `json:"address" gorm:"column:address;type:varchar(255);not null"`
	TotalFloor int    `json:"totalFloor" gorm:"column:total_floor;type:int;not null"`
	TotalRoom  int    `json:"totalRoom" gorm:"column:total_room;type:int;not null"`
}

func (u *BuildingModel) TableName() string {
	return "building"
}

func (u *BuildingModel) BeforeCreate(tx *gorm.DB) error {
	username, _ := tx.Get("username")
	if username == nil {
		username = "SYSTEM"
	}
	u.CreatedBy = username.(string)
	u.UpdatedBy = username.(string)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return nil
}

func (u *BuildingModel) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("updated_at", "updated_by") {
		return errors.New(config.GetMessageCode("CONCURRENCY_ERROR"))
	}

	username, _ := tx.Get("username")
	if username == nil {
		username = "SYSTEM"
	}
	tx.Statement.SetColumn("created_by", username.(string))
	tx.Statement.SetColumn("updated_by", username.(string))
	tx.Statement.SetColumn("created_at", time.Now())
	tx.Statement.SetColumn("updated_at", time.Now())

	return nil
}
