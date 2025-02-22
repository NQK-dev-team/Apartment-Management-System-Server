package migrations

import (
	"api/config"
	"api/models"
)

type EmailQueueMigration struct{}

func NewEmailQueueMigration() *EmailQueueMigration {
	return &EmailQueueMigration{}
}

func (m *EmailQueueMigration) Up() {
	config.DB.AutoMigrate(&models.EmailQueueModel{})
	config.DB.AutoMigrate(&models.EmailQueueFailModel{})
}

func (m *EmailQueueMigration) Down() {
	config.DB.Migrator().DropTable(&models.EmailQueueModel{})
	config.DB.Migrator().DropTable(&models.EmailQueueFailModel{})
}
