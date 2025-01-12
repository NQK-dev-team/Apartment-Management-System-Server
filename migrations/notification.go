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
	config.DB.AutoMigrate(&models.NotificationReceiverModel{})
}

func (m *NotificationMigration) Down() {
	config.DB.Migrator().DropTable(&models.NotificationReceiverModel{})
	config.DB.Migrator().DropTable(&models.NotificationFileModel{})
	config.DB.Migrator().DropTable(&models.NotificationModel{})
}
