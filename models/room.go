package models

import (
	"time"

	"gorm.io/gorm"
)

type RoomModel struct {
	DefaultModel
	No          int              `json:"no" gorm:"column:no;type:int;not null;"`
	Floor       int              `json:"floor" gorm:"column:floor;type:int;not null;"`
	Description string           `json:"description" gorm:"column:description;type:varchar(255);"`
	Area        float64          `json:"area" gorm:"column:area;type:numeric;not null;"`
	Status      int              `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Rented, 2: Bought, 3: Available, 4: Maintenanced, 5: Unavailable
	BuildingID  int64            `json:"buildingID" gorm:"column:building_id;not null;"`
	Contracts   []ContractModel  `json:"contracts" gorm:"foreignKey:room_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Images      []RoomImageModel `json:"images" gorm:"foreignKey:room_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Contracts   []ContractModel  `json:"contracts" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Images      []RoomImageModel `json:"images" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Building    BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *RoomModel) TableName() string {
	return "room"
}

// func (u *RoomModel) BeforeDelete(tx *gorm.DB) error {
// 	userID, _ := tx.Get("userID")

// 	return tx.Transaction(func(tx1 *gorm.DB) error {
// 		if err := tx1.Set("userID", userID).Where("room_id = ?", u.ID).Delete(&ContractModel{}).Error; err != nil {
// 			return err
// 		}

// 		if err := tx1.Set("userID", userID).Where("room_id = ?", u.ID).Delete(&RoomImageModel{}).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// func (u *RoomModel) BeforeCreate(tx *gorm.DB) error {
// 	lastRoom := RoomModel{}
// 	// Get the last room of the building floor
// 	tx.Where("building_id = ? AND floor = ?", u.BuildingID, u.Floor).Order("no desc").First(&lastRoom)

// 	if lastRoom.No == 0 {
// 		u.No = 1000*u.Floor + 1
// 	} else {
// 		u.No = lastRoom.No + 1
// 	}

// 	return nil
// }

func (u *RoomModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *RoomModel) BeforeUpdate(tx *gorm.DB) error {
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
