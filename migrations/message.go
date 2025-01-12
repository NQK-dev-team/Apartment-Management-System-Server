package migrations

import (
	"api/config"
	"api/models"
)

type MessageMigration struct {
}

func NewMessageMigration() *MessageMigration {
	return &MessageMigration{}
}

func (m *MessageMigration) Up() {
	config.DB.AutoMigrate(&models.MessageModel{})
	config.DB.AutoMigrate(&models.MessageFileModel{})
}

func (m *MessageMigration) Down() {
	config.DB.Migrator().DropTable(&models.MessageModel{})
	config.DB.Migrator().DropTable(&models.MessageFileModel{})
}
