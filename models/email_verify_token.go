package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerifyTokenModel struct {
	ID        int64     `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	Token     string    `json:"token" gorm:"column:token;type:varchar(255);not null;"`
	Email     string    `json:"email" gorm:"column:email;type:varchar(255);not null;"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expires_at;type:timestamp with time zone;not null;"`
}

func (u *EmailVerifyTokenModel) TableName() string {
	return "email_verify_token"
}

func (u *EmailVerifyTokenModel) BeforeCreate(tx *gorm.DB) error {
	// u.CreatedAt = time.Now()
	u.ExpiresAt = time.Now().Add(24 * 7 * time.Hour)

	return nil
}
