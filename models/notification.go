package models

import "time"

type NotificationModel struct {
	DefaultModel
	Title    string    `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content  string    `json:"content" gorm:"column:content;type:text;not null;"`
	SendTime time.Time `json:"send_time" gorm:"column:send_time;type:timestamp with time zone;not null;default:now();"`
	SenderID int64     `json:"sender_id" gorm:"column:sender_id;not null;"`
	Sender   UserModel `json:"sender" gorm:"foreignKey:sender_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *NotificationModel) TableName() string {
	return "notification"
}
