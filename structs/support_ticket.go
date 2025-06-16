package structs

import (
	"api/models"
	"database/sql"
)

type SupportTicket struct {
	models.DefaultModel
	Status             int                             `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Pending, 2: Approved, 3: Rejected
	Title              string                          `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content            string                          `json:"content" gorm:"column:content;type:text;not null;"`
	ContractID         int64                           `json:"contractID" gorm:"column:contract_id;not null;"`
	Contract           models.ContractModel            `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID         int64                           `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer           models.UserModel                `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Files              []models.SupportTicketFileModel `json:"files" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerID          int64                           `json:"managerID" gorm:"column:manager_id;"`
	Manager            models.UserModel                `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerResult      bool                            `json:"managerResult" gorm:"column:manager_result;type:bool;"` // 0: Rejected, 1: Approved
	ManagerResolveTime sql.NullTime                    `json:"managerResolveTime" gorm:"column:manager_resolve_time;type:timestamp with time zone;"`
	OwnerID            int64                           `json:"ownerID" gorm:"column:owner_id;"`
	Owner              models.UserModel                `json:"owner" gorm:"foreignKey:owner_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OwnerResult        bool                            `json:"ownerResult" gorm:"column:owner_result;type:bool;"` // 0: Rejected, 1: Approved
	OwnerResolveTime   sql.NullTime                    `json:"ownerResolveTime" gorm:"column:owner_resolve_time;type:timestamp with time zone;"`
	BuildingName       string                          `json:"buildingName" gorm:"building_name"`
	RoomNo             int                             `json:"roomNo" gorm:"room_no"`
	RoomFloor          int                             `json:"roomFloor" gorm:"room_floor"`
}
