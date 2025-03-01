package models

type SupportTicketFileModel struct {
	DefaultFileModel
	SupportTicketID int64 `json:"supportTicketID" gorm:"column:support_ticket_id;not null;"`
	// SupportTicket   SupportTicketModel `json:"supportTicket" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *SupportTicketFileModel) TableName() string {
	return "support_ticket_file"
}
