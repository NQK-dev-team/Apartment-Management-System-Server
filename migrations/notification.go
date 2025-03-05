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
	config.DB.AutoMigrate(&models.NotificationModel{})
	config.DB.AutoMigrate(&models.NotificationFileModel{})
	config.DB.Exec("ALTER TABLE notification_file ADD CONSTRAINT notification_file_composite_key UNIQUE (notification_id, id);")
	config.DB.AutoMigrate(&models.NotificationReceiverModel{})
	config.DB.Exec("ALTER TABLE notification_receiver ADD CONSTRAINT notification_receiver_composite_key UNIQUE (notification_id, user_id);")
}

func (m *NotificationMigration) Down() {
	config.DB.Migrator().DropTable(&models.NotificationReceiverModel{})
	config.DB.Migrator().DropTable(&models.NotificationFileModel{})
	config.DB.Migrator().DropTable(&models.NotificationModel{})
}
