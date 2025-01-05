package migrations

import (
	"api/config"
	"api/models"
)

type RefreshTokenMigration struct {
}

func NewRefreshTokenMigration() *RefreshTokenMigration {
	return &RefreshTokenMigration{}
}

func (m *RefreshTokenMigration) Up() {
	config.DB.AutoMigrate(&models.RefreshTokenModel{})
}

func (m *RefreshTokenMigration) Down() {
	config.DB.Migrator().DropTable(&models.RefreshTokenModel{})
}
