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
	// Add role check constraint
	config.DB.Exec("ALTER TABLE \"user\" ADD CONSTRAINT role_check CHECK ((NOT (is_customer AND (is_owner OR is_manager))))")
}

func (m *UserMigration) Down() {
	config.DB.Migrator().DropTable(&models.UserModel{})
}
