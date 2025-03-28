package structs

import (
	"api/models"
	"time"
)

type SupportTicket struct {
	models.DefaultModel
	Status       int                             `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Pending, 2: Approved, 3: Rejected
	Title        string                          `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content      string                          `json:"content" gorm:"column:content;type:text;not null;"`
	ContractID   int64                           `json:"contractID" gorm:"column:contract_id;not null;"`
	Contract     models.ContractModel            `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID   int64                           `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer     models.UserModel                `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Files        []models.SupportTicketFileModel `json:"files" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomNo       int                             `json:"roomnNo" gorm:"column:room_no;type:int;not null;"`
	BuildingName string                          `json:"buildingName" gorm:"column:building_name;type:varchar(255);not null;"`
}

func (u *SupportTicket) TableName() string {
	return "support_ticket"
}

type ResolveTicket struct {
	SupportTicketID int64            `json:"supportTicketID" gorm:"column:support_ticket_id;not null;"`
	SupportTicket   SupportTicket    `json:"supportTicket" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerID       int64            `json:"managerID" gorm:"column:manager_id;not null;"`
	Manager         models.UserModel `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Result          bool             `json:"result" gorm:"column:result;type:bool;not null;"` // 0: Rejected, 1: Approved
	ResolveTime     time.Time        `json:"resolveTime" gorm:"column:resolve_time;type:timestamp with time zone;not null;default:now()"`
	OwnerID         int64            `json:"ownerID"`
	Owner           models.UserModel `json:"owner"`
}
