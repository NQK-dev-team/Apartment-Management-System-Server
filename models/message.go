package models

import "time"

type MessageModel struct {
	DefaultModel
	Content    string    `json:"content" gorm:"column:content;type:text;not null;"`
	SendAt     time.Time `json:"sendAt" gorm:"column:send_at;type:timestamp with time zone;not null;default:now();"`
	Sender     int       `json:"sender" gorm:"column:sender;not null;type:int;"` // 0: Manager, 1: Customer
	ManagerID  int64     `json:"managerID" gorm:"column:manager_id;not null;"`
	Manager    UserModel `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CustomerID int64     `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer   UserModel `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *MessageModel) TableName() string {
	return "message"
}
