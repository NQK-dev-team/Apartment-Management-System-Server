package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageModel struct {
	DefaultModel
	Content    string             `json:"content" gorm:"column:content;type:text;not null;"`
	SendAt     time.Time          `json:"sendAt" gorm:"column:send_at;type:timestamp with time zone;not null;default:now();"`
	Sender     int                `json:"sender" gorm:"column:sender;not null;type:int;"` // 0: Manager, 1: Customer
	ManagerID  int64              `json:"managerID" gorm:"column:manager_id;not null;"`
	Manager    UserModel          `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CustomerID int64              `json:"customerID" gorm:"column:customer_id;not null;"`
	Customer   UserModel          `json:"customer" gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Files      []MessageFileModel `json:"files" gorm:"foreignKey:message_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *MessageModel) TableName() string {
	return "message"
}

// func (u *MessageModel) BeforeDelete(tx *gorm.DB) error {
// 	userID, _ := tx.Get("userID")

// 	return tx.Transaction(func(tx1 *gorm.DB) error {
// 		if err := tx1.Set("userID", userID).Model(&MessageFileModel{}).Where("message_id = ?", u.ID).Delete(&MessageFileModel{}).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

func (u *MessageModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *MessageModel) BeforeUpdate(tx *gorm.DB) error {
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
