package migrations

import (
	"api/config"
	"api/models"
)

type EmailVerifyTokenMigration struct {
}

func NewEmailVerifyTokenMigration() *EmailVerifyTokenMigration {
	return &EmailVerifyTokenMigration{}
}

func (m *EmailVerifyTokenMigration) Up() {
	config.DB.AutoMigrate(&models.EmailVerifyTokenModel{})
}

func (m *EmailVerifyTokenMigration) Down() {
	config.DB.Migrator().DropTable(&models.EmailVerifyTokenModel{})
}
