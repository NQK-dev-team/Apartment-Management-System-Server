package models

import (
	"api/config"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ManagerResolveSupportTicketModel struct {
	SupportTicketID int64              `json:"supportTicketID" gorm:"column:support_ticket_id;primaryKey;"`
	SupportTicket   SupportTicketModel `json:"supportTicket" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerID       int64              `json:"managerID" gorm:"column:manager_id;primaryKey;"`
	Manager         UserModel          `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Result          int                `json:"result" gorm:"column:result;type:int;not null;"` // 0: Rejected, 1: Approved
	ResolveTime     time.Time          `json:"resolveTime" gorm:"column:resolve_time;type:timestamp with time zone;not null;default:now()"`
}

func (u *ManagerResolveSupportTicketModel) TableName() string {
	return "manager_resolve_support_ticket"
}

func (u *ManagerResolveSupportTicketModel) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("updated_at", "updated_by") {
		return errors.New(config.GetMessageCode("CONCURRENCY_ERROR"))
	}

	userID, _ := tx.Get("userID")
	if userID == nil {
		userID = "SYSTEM"
	}
	tx.Statement.SetColumn("updated_by", userID.(string))
	tx.Statement.SetColumn("updated_at", time.Now())

	return nil
}
