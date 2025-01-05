package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerifyTokenModel struct {
	ID        int64     `json:"id" gorm:"primaryKey; colunn:id; autoIncrement; not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null"`
	Token     string    `json:"token" gorm:"column:token;type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"column:email;type:varchar(255);not null"`
}

func (a *EmailVerifyTokenModel) TableName() string {
	return "email_verify_token"
}

func (a *EmailVerifyTokenModel) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()

	return nil
}
