package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type SupportTicketModel struct {
	DefaultModel
	Status             int                      `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Pending, 2: Approved, 3: Rejected
	Title              string                   `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content            string                   `json:"content" gorm:"column:content;type:text;not null;"`
	ContractID         int64                    `json:"contractID" gorm:"column:contract_id;not null;"`
	Contract           ContractModel            `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID         int64                    `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer           UserModel                `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Files              []SupportTicketFileModel `json:"files" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerID          int64                    `json:"managerID" gorm:"column:manager_id;"`
	Manager            UserModel                `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManagerResult      sql.NullBool             `json:"managerResult" gorm:"column:manager_result;type:bool;"` // 0: Rejected, 1: Approved
	ManagerResolveTime sql.NullTime             `json:"managerResolveTime" gorm:"column:manager_resolve_time;type:timestamp with time zone;"`
	OwnerID            int64                    `json:"ownerID" gorm:"column:owner_id;"`
	Owner              UserModel                `json:"owner" gorm:"foreignKey:owner_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OwnerResult        sql.NullBool             `json:"ownerResult" gorm:"column:owner_result;type:bool;"` // 0: Rejected, 1: Approved
	OwnerResolveTime   sql.NullTime             `json:"ownerResolveTime" gorm:"column:owner_resolve_time;type:timestamp with time zone;"`
}

func (u *SupportTicketModel) TableName() string {
	return "support_ticket"
}

func (u *SupportTicketModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *SupportTicketModel) BeforeUpdate(tx *gorm.DB) error {
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
