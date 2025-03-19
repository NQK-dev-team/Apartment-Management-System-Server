package models

import (
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

func (u *BuildingModel) BeforeDelete(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")

	return tx.Transaction(func(tx1 *gorm.DB) error {
		if err := tx1.Set("userID", userID).Where("building_id = ?", u.ID).Delete(&BuildingImageModel{}).Error; err != nil {
			return err
		}

		if err := tx1.Set("userID", userID).Where("building_id = ?", u.ID).Delete(&BuildingServiceModel{}).Error; err != nil {
			return err
		}

		if err := tx1.Set("userID", userID).Where("building_id = ?", u.ID).Delete(&RoomModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}
