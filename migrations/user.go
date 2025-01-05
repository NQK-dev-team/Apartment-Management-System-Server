package migrations

import (
	"api/config"
	"api/models"
)

type UserMigration struct {
}

func NewUserMigration() *UserMigration {
	return &UserMigration{}
}

func (m *UserMigration) Up() {
	config.DB.AutoMigrate(&models.UserModel{})
}

func (m *UserMigration) Down() {
	config.DB.Migrator().DropTable(&models.UserModel{})
}
