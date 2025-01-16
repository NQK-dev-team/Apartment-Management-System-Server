package models

type SupportTicketModel struct {
	DefaultModel
	Status     int           `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Pending, 2: Approved, 3: Rejected
	Title      string        `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content    string        `json:"content" gorm:"column:content;type:text;not null;"`
	ContractID int64         `json:"contractID" gorm:"column:contract_id;not null;"`
	Contract   ContractModel `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID int64         `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer   UserModel     `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *SupportTicketModel) TableName() string {
	return "support_ticket"
}
