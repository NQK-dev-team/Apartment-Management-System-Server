package models

import (
	"api/config"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenModel struct {
	ID        int64     `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	Token     string    `json:"token" gorm:"column:token;type:varchar(255);not null;"`
	UserID    int64     `json:"userID" gorm:"column:user_id;not null;"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expires_at;type:timestamp with time zone;not null;"`
}

func (u *RefreshTokenModel) TableName() string {
	return "refresh_token"
}

func (u *RefreshTokenModel) BeforeCreate(tx *gorm.DB) error {
	// u.CreatedAt = time.Now()
	expirationTimeStr, err := config.GetEnv("JWT_REFRESH_EXPIRE_TIME")
	if err != nil {
		expirationTimeStr = "604800"
	}
	expirationTime, err := strconv.Atoi(expirationTimeStr)
	if err != nil {
		expirationTime = 604800
	}
	u.ExpiresAt = time.Now().Add(time.Second * time.Duration(expirationTime))

	return nil
}
