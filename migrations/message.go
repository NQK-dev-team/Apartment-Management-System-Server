package migrations

// This file is not used in the current implementation.

// import (
// 	"api/config"
// 	"api/models"
// )

// type MessageMigration struct {
// }

// func NewMessageMigration() *MessageMigration {
// 	return &MessageMigration{}
// }

// func (m *MessageMigration) Up() {
// 	config.MigrationDB.AutoMigrate(&models.MessageModel{})
// 	config.DB.Exec("ALTER TABLE message ADD CONSTRAINT message_sender CHECK (sender>=0 and sender<=1);")
// 	config.MigrationDB.AutoMigrate(&models.MessageFileModel{})
// 	// config.DB.Exec("ALTER TABLE message_file ADD CONSTRAINT message_file_composite_key UNIQUE (message_id, id);")
// }

// func (m *MessageMigration) Down() {
// 	config.MigrationDB.Migrator().DropTable(&models.MessageModel{})
// 	config.MigrationDB.Migrator().DropTable(&models.MessageFileModel{})
// }
