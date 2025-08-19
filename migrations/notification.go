package migrations

import (
	"api/config"
	"api/models"
)

type NotificationMigration struct{}

func NewNotificationMigration() *NotificationMigration {
	return &NotificationMigration{}
}

func (m *NotificationMigration) Up() {
	config.MigrationDB.AutoMigrate(&models.NotificationModel{})
	config.MigrationDB.AutoMigrate(&models.NotificationFileModel{})
	// config.DB.Exec("ALTER TABLE notification_file ADD CONSTRAINT notification_file_composite_key UNIQUE (notification_id, id);")
	config.MigrationDB.AutoMigrate(&models.NotificationReceiverModel{})
	// config.DB.Exec("ALTER TABLE notification_receiver ADD CONSTRAINT notification_receiver_composite_key UNIQUE (notification_id, user_id);")
}

func (m *NotificationMigration) Down() {
	config.MigrationDB.Migrator().DropTable(&models.NotificationReceiverModel{})
	config.MigrationDB.Migrator().DropTable(&models.NotificationFileModel{})
	config.MigrationDB.Migrator().DropTable(&models.NotificationModel{})
}
