package structs

import (
	"api/models"
	"time"
)

type SupportTicket struct {
	ID              int64                           `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	Status          int                             `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Pending, 2: Approved, 3: Rejected
	Title           string                          `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content         string                          `json:"content" gorm:"column:content;type:text;not null;"`
	ContractID      int64                           `json:"contractID" gorm:"column:contract_id;not null;"`
	CustomerID      int64                           `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer        models.UserModel                `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Files           []models.SupportTicketFileModel `json:"files" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt       time.Time                       `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	SupportTicketID int64                           `json:"supportTicketID" gorm:"column:support_ticket_id;not null;"`
	SupportTicket   models.SupportTicketModel       `json:"supportTicket" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerID       int64                           `json:"managerID" gorm:"column:manager_id;not null;"`
	Manager         models.UserModel                `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Result          bool                            `json:"result" gorm:"column:result;type:bool;not null;"` // 0: Rejected, 1: Approved
}
