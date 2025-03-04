package models

import "gorm.io/gorm"

type SupportTicketModel struct {
	DefaultModel
	Status     int                      `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Pending, 2: Approved, 3: Rejected
	Title      string                   `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content    string                   `json:"content" gorm:"column:content;type:text;not null;"`
	ContractID int64                    `json:"contractID" gorm:"column:contract_id;not null;"`
	Contract   ContractModel            `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID int64                    `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer   UserModel                `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Files      []SupportTicketFileModel `json:"files" gorm:"foreignKey:support_ticket_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *SupportTicketModel) TableName() string {
	return "support_ticket"
}

func (u *SupportTicketModel) BeforeDelete(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")

	return tx.Transaction(func(tx1 *gorm.DB) error {
		if err := tx1.Set("userID", userID).Model(&SupportTicketFileModel{}).Where("support_ticket_id = ?", u.ID).Delete(&SupportTicketFileModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}
