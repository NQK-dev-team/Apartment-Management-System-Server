package models

import (
	"api/config"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenModel struct {
	ID        int64     `json:"id" gorm:"primaryKey; colunn:id; autoIncrement; not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null"`
	Token     string    `json:"token" gorm:"column:token;type:varchar(255);not null"`
	UserID    int64     `json:"userId" gorm:"column:user_id;not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expires_at;type:timestamp with time zone;not null"`
}

func (a *RefreshTokenModel) TableName() string {
	return "refresh_token"
}

func (a *RefreshTokenModel) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	expirationTimeStr, err := config.GetEnv("JWT_REFRESH_EXPIRE_TIME")
	if err != nil {
		expirationTimeStr = "604800"
	}
	expirationTime, err := strconv.Atoi(expirationTimeStr)
	if err != nil {
		expirationTime = 604800
	}
	a.ExpiresAt = time.Now().Add(time.Second * time.Duration(expirationTime))

	return nil
}
