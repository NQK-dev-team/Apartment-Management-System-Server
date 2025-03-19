package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationModel struct {
	DefaultModel
	Title    string                  `json:"title" gorm:"column:title;type:varchar(255);not null;"`
	Content  string                  `json:"content" gorm:"column:content;type:text;not null;"`
	SendTime time.Time               `json:"sendTime" gorm:"column:send_time;type:timestamp with time zone;not null;default:now();"`
	SenderID int64                   `json:"senderID" gorm:"column:sender_id;not null;"`
	Sender   UserModel               `json:"sender" gorm:"foreignKey:sender_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Files    []NotificationFileModel `json:"files" gorm:"foreignKey:notification_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *NotificationModel) TableName() string {
	return "notification"
}

func (u *NotificationModel) BeforeDelete(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")

	return tx.Transaction(func(tx1 *gorm.DB) error {
		if err := tx1.Set("userID", userID).Model(&NotificationFileModel{}).Where("notification_id = ?", u.ID).Delete(&NotificationFileModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}
