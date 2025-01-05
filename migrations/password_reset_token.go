package migrations

import (
	"api/config"
	"api/models"
)

type PasswordResetTokenMigration struct {
}

func NewPasswordResetTokenMigration() *PasswordResetTokenMigration {
	return &PasswordResetTokenMigration{}
}

func (m *PasswordResetTokenMigration) Up() {
	config.DB.AutoMigrate(&models.PasswordResetTokenModel{})
}

func (m *PasswordResetTokenMigration) Down() {
	config.DB.Migrator().DropTable(&models.PasswordResetTokenModel{})
}
