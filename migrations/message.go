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
	config.DB.Exec("ALTER TABLE message ADD CONSTRAINT message_sender CHECK (sender>=0 and sender<=1);")
	config.DB.AutoMigrate(&models.MessageFileModel{})
	config.DB.Exec("ALTER TABLE message_file ADD CONSTRAINT message_file_composite_key UNIQUE (message_id, id);")
}

func (m *MessageMigration) Down() {
	config.DB.Migrator().DropTable(&models.MessageModel{})
	config.DB.Migrator().DropTable(&models.MessageFileModel{})
}
