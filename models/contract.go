package models

import (
	"database/sql"
	"time"
)

type ContractModel struct {
	DefaultModel
	Status         int                  `json:"status" gorm:"column:status;type:int;not null;"` // 1: Active, 2: Expired, 3: Cancelled, 4: Waiting for signatures, 5: Not in effect yet
	Value          float64              `json:"value" gorm:"column:value;type:float;not null;"`
	Type           int                  `json:"type" gorm:"column:type;type:int;not null;"` // 1: Rent, 2: Buy
	StartDate      time.Time            `json:"startDate" gorm:"column:start_date;type:date;not null;default:now();"`
	EndDate        sql.NullTime         `json:"endDate" gorm:"column:end_date;type:date;"`
	SignDate       sql.NullTime         `json:"signDate" gorm:"column:sign_date;type:date;default:now();"`
	CreatorID      int64                `json:"creatorID" gorm:"column:creator_id;not null;"`
	Creator        UserModel            `json:"creator" gorm:"foreignKey:creator_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HouseholderID  int64                `json:"householderID" gorm:"column:householder_id;not null;"`
	Householder    UserModel            `json:"householder" gorm:"foreignKey:householder_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomID         int64                `json:"roomID" gorm:"column:room_id;not null;"`
	Bills          []BillModel          `json:"bills" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Files          []ContractFileModel  `json:"files" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SupportTickets []SupportTicketModel `json:"supportTickets" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// BuildingID    int64     `json:"buildingID" gorm:"column:building_id;not null;"`
	// Room          RoomModel           `json:"room" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *ContractModel) TableName() string {
	return "contract"
}

// func (u *ContractModel) BeforeDelete(tx *gorm.DB) (err error) {
// 	userID, _ := tx.Get("userID")

// 	return tx.Transaction(func(tx1 *gorm.DB) error {
// 		if err := tx1.Set("userID", userID).Model(&ContractFileModel{}).Where("contract_id = ?", u.ID).Delete(&ContractFileModel{}).Error; err != nil {
// 			return err
// 		}

// 		if err := tx1.Set("userID", userID).Model(&BillModel{}).Where("contract_id = ?", u.ID).Delete(&BillModel{}).Error; err != nil {
// 			return err
// 		}

// 		if err := tx1.Set("userID", userID).Model(&RoomResidentModel{}).Model(&RoomResidentListModel{}).
// 			Joins("JOIN room_resident_list ON room_resident.id = room_resident_list.resident_id").
// 			Where("room_resident_list.contract_id = ?", u.ID).
// 			Delete(&RoomResidentModel{}).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }
