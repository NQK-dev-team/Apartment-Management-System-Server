package models

import "time"

type ContractModel struct {
	DefaultModel
	Status        int       `json:"status" gorm:"column:status;type:int;not null;"` // 0: Active, 1: Expired, 2: Cancelled, 3: Waiting for signatures, 4: Not in effect yet
	Value         float64   `json:"value" gorm:"column:value;type:float;not null;"`
	Type          int       `json:"type" gorm:"column:type;type:int;not null;"` // 0: Rent, 1: Buy
	StartDate     time.Time `json:"startDate" gorm:"column:start_date;type:date;not null;default:now();"`
	EndDate       time.Time `json:"endDate" gorm:"column:end_date;type:date;default:now();"`
	SignDate      time.Time `json:"signDate" gorm:"column:sign_date;type:date;default:now();"`
	CreatorID     int64     `json:"creatorID" gorm:"column:creator_id;not null;"`
	Creator       UserModel `json:"creator" gorm:"foreignKey:creator_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HouseholderID int64     `json:"householderID" gorm:"column:householder_id;not null;"`
	Householder   UserModel `json:"householder" gorm:"foreignKey:householder_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomID        int64     `json:"roomID" gorm:"column:room_id;not null;"`
	BuildingID    int64     `json:"buildingID" gorm:"column:building_id;not null;"`
	Room          RoomModel `json:"room" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *ContractModel) TableName() string {
	return "contract"
}
