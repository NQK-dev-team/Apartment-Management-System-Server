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

// func (u *NotificationModel) BeforeDelete(tx *gorm.DB) error {
// 	userID, _ := tx.Get("userID")

// 	return tx.Transaction(func(tx1 *gorm.DB) error {
// 		if err := tx1.Set("userID", userID).Model(&NotificationFileModel{}).Where("notification_id = ?", u.ID).Delete(&NotificationFileModel{}).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

func (u *NotificationModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *NotificationModel) BeforeUpdate(tx *gorm.DB) error {
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
