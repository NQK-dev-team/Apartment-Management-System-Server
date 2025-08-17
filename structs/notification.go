package structs

import (
	"api/models"
	"mime/multipart"
	"time"
)

type NewNotification struct {
	Title       string  `form:"title" validate:"required"`
	Content     string  `form:"content" validate:"required"`
	ReceiverStr string  `form:"receiverStr"`
	Receivers   []int64 `form:"-" validate:"required"`
	Files       []*multipart.FileHeader
}

type Notification struct {
	models.DefaultModel
	Title     string                         `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content   string                         `json:"content" gorm:"column:content;type:text;not null;"`
	SendTime  time.Time                      `json:"sendTime" gorm:"column:send_time;type:timestamp with time zone;not null;default:now();"`
	SenderID  int64                          `json:"senderID" gorm:"column:sender_id;not null;"`
	Sender    models.UserModel               `json:"sender" gorm:"foreignKey:sender_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Files     []models.NotificationFileModel `json:"files" gorm:"foreignKey:notification_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsRead    int                            `json:"isRead" gorm:"column:is_read;"`
	IsMarked  int                            `json:"isMarked" gorm:"column:is_marked;"`
	Receivers []int64                        `json:"receivers" gorm:"-"`
}

// type NotificationDetail struct {
// 	models.DefaultModel
// 	Title     string                         `json:"title" gorm:"column:title;type:varchar(255);not null;"`
// 	Content   string                         `json:"content" gorm:"column:content;type:text;not null;"`
// 	SendTime  time.Time                      `json:"sendTime" gorm:"column:send_time;type:timestamp with time zone;not null;default:now();"`
// 	SenderID  int64                          `json:"senderID" gorm:"column:sender_id;not null;"`
// 	Sender    models.UserModel               `json:"sender" gorm:"foreignKey:sender_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
// 	Files     []models.NotificationFileModel `json:"files" gorm:"foreignKey:notification_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// 	Receivers []int64                        `json:"receivers" gorm:"-"`
// }