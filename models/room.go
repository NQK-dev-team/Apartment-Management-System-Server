package models

import (
	"api/config"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RoomModel struct {
	DefaultModel
	No          int           `json:"no" gorm:"column:no;type:int;not null"`
	Floor       int           `json:"floor" gorm:"column:floor;type:int;not null"`
	Description string        `json:"description" gorm:"column:description;type:varchar(255)"`
	Area        float32       `json:"area" gorm:"column:area;type:numeric;not null"`
	Status      int           `json:"status" gorm:"column:status;type:int;not null;default:1"`
	BuildingID  int64         `json:"buildingId" gorm:"column:building_id;not null"`
	Building    BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id"`
}

func (u *RoomModel) TableName() string {
	return "room"
}

func (u *RoomModel) BeforeCreate(tx *gorm.DB) error {
	username, _ := tx.Get("username")
	if username == nil {
		username = "SYSTEM"
	}
	u.CreatedBy = username.(string)
	u.UpdatedBy = username.(string)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	lastRoom := RoomModel{}
	// Get the last room of the building floor
	tx.Where("building_id = ? AND floor = ?", u.BuildingID, u.Floor).Order("no desc").First(&lastRoom)

	if lastRoom.No == 0 {
		u.No = 1000*u.Floor + 1
	} else {
		u.No = lastRoom.No + 1
	}

	return nil
}

func (u *RoomModel) BeforeUpdate(tx *gorm.DB) error {
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
