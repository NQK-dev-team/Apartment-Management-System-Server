package models

import (
	"gorm.io/gorm"
)

type RoomModel struct {
	DefaultModel
	No          int           `json:"no" gorm:"column:no;type:int;not null;"`
	Floor       int           `json:"floor" gorm:"column:floor;type:int;not null;"`
	Description string        `json:"description" gorm:"column:description;type:varchar(255);"`
	Area        float32       `json:"area" gorm:"column:area;type:numeric;not null;"`
	Status      int           `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Rented, 2: Bought, 3: Available, 4: Maintenanced, 5: Unavailable
	BuildingID  int64         `json:"buildingID" gorm:"column:building_id;primaryKey;"`
	Building    BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *RoomModel) TableName() string {
	return "room"
}

func (u *RoomModel) BeforeCreate(tx *gorm.DB) error {
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
