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
	config.DB.Exec("ALTER TABLE support_ticket ADD CONSTRAINT support_ticket_status CHECK (status >= 1 AND status <= 3);")
	config.DB.AutoMigrate(&models.SupportTicketFileModel{})
	config.DB.Exec("ALTER TABLE support_ticket_file ADD CONSTRAINT support_ticket_file_composite_key UNIQUE (id, support_ticket_id);")
	config.DB.AutoMigrate(&models.ManagerResolveSupportTicketModel{})
	config.DB.Exec("ALTER TABLE manager_resolve_support_ticket ADD CONSTRAINT manager_resolve_support_ticket_composite_key UNIQUE (manager_id, support_ticket_id);")
}

func (m *SupportTicketMigration) Down() {
	config.DB.Migrator().DropTable(&models.ManagerResolveSupportTicketModel{})
	config.DB.Migrator().DropTable(&models.SupportTicketFileModel{})
	config.DB.Migrator().DropTable(&models.SupportTicketModel{})
}
