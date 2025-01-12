package models

type RoomResidentListModel struct {
	ContractID int64             `json:"contractID" gorm:"column:contract_id;type:int;primaryKey;"`
	ResidentID int64             `json:"residentID" gorm:"column:resident_id;type:int;primaryKey;"`
	Contract   ContractModel     `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Resident   RoomResidentModel `json:"resident" gorm:"foreignKey:resident_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *RoomResidentListModel) TableName() string {
	return "room_resident_list"
}
