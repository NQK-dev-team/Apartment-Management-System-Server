package migrations

import (
	"api/config"
	"api/models"
)

type SystemMigration struct{}

func NewSystemMigration() *SystemMigration {
	return &SystemMigration{}
}

func (m *SystemMigration) Up() {
	config.MigrationDB.AutoMigrate(&models.EmailQueueModel{})
	config.MigrationDB.AutoMigrate(&models.EmailQueueFailModel{})
	config.MigrationDB.AutoMigrate(&models.EmailVerifyTokenModel{})
	config.MigrationDB.AutoMigrate(&models.RefreshTokenModel{})
	config.MigrationDB.AutoMigrate(&models.PasswordResetTokenModel{})
	config.MigrationDB.AutoMigrate(&models.UploadFileModel{})
}

func (m *SystemMigration) Down() {
	config.MigrationDB.Migrator().DropTable(&models.EmailQueueModel{})
	config.MigrationDB.Migrator().DropTable(&models.EmailQueueFailModel{})
	config.MigrationDB.Migrator().DropTable(&models.EmailVerifyTokenModel{})
	config.MigrationDB.Migrator().DropTable(&models.RefreshTokenModel{})
	config.MigrationDB.Migrator().DropTable(&models.PasswordResetTokenModel{})
	config.MigrationDB.Migrator().DropTable(&models.UploadFileModel{})
}
