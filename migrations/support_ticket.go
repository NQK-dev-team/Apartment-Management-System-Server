package migrations

import (
	"api/config"
	"api/models"
)

type SupportTicketMigration struct {
}

func NewSupportTicketMigration() *SupportTicketMigration {
	return &SupportTicketMigration{}
}

func (m *SupportTicketMigration) Up() {
	config.DB.AutoMigrate(&models.SupportTicketModel{})
	config.DB.AutoMigrate(&models.SupportTicketFileModel{})
	config.DB.AutoMigrate(&models.ManagerResolveSupportTicketModel{})
	config.DB.Exec("ALTER TABLE support_ticket ADD CONSTRAINT support_ticket_status CHECK (status >= 1 AND status <= 3);")
}

func (m *SupportTicketMigration) Down() {
	config.DB.Migrator().DropTable(&models.ManagerResolveSupportTicketModel{})
	config.DB.Migrator().DropTable(&models.SupportTicketFileModel{})
	config.DB.Migrator().DropTable(&models.SupportTicketModel{})
}
