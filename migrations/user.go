package migrations

import (
	"api/config"
	"api/models"
	"fmt"
)

type UserMigration struct {
}

func NewUserMigration() *UserMigration {
	return &UserMigration{}
}

func (m *UserMigration) generateRoleCheckRule() string {
	adminRole := "(is_owner AND is_manager AND NOT is_customer)"
	managerRole := "(NOT is_owner AND is_manager AND NOT is_customer)"
	customerRole := "(NOT is_owner AND NOT is_manager AND is_customer)"
	return fmt.Sprintf("ALTER TABLE \"user\" ADD CONSTRAINT role_check CHECK (%s OR %s OR %s)", adminRole, managerRole, customerRole)
}

func (m *UserMigration) Up() {
	config.DB.AutoMigrate(&models.UserModel{})
	// Add role check constraint
	config.DB.Exec(m.generateRoleCheckRule())
}

func (m *UserMigration) Down() {
	config.DB.Migrator().DropTable(&models.UserModel{})
}
