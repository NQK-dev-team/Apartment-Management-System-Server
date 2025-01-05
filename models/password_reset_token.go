package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetTokenModel struct {
	ID        int64     `json:"id" gorm:"primaryKey; colunn:id; autoIncrement; not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null"`
	Token     string    `json:"token" gorm:"column:token;type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"column:email;type:varchar(255);not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expires_at;type:timestamp with time zone;not null"`
}

func (a *PasswordResetTokenModel) TableName() string {
	return "password_reset_token"
}

func (a *PasswordResetTokenModel) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()

	return nil
}
