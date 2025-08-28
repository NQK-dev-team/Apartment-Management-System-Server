package models

import (
	"time"

	"gorm.io/gorm"
)

type SupportTicketFileModel struct {
	DefaultFileModel
	SupportTicketID int64 `json:"supportTicketID" gorm:"column:support_ticket_id;not null;"`
	// SupportTicket   SupportTicketModel `json:"supportTicket" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *SupportTicketFileModel) TableName() string {
	return "support_ticket_file"
}

func (u *SupportTicketFileModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
	}
	u.CreatedAt = time.Now()

	return nil
}
