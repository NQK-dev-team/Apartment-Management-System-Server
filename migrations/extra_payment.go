package migrations

import (
	"api/config"
	"api/models"
)

type ExtraPaymentMigration struct {
}

func NewExtraPaymentMigration() *ExtraPaymentMigration {
	return &ExtraPaymentMigration{}
}

func (m *ExtraPaymentMigration) Up() {
	config.DB.AutoMigrate(&models.ExtraPaymentModel{})
}

func (m *ExtraPaymentMigration) Down() {
	config.DB.Migrator().DropTable(&models.ExtraPaymentModel{})
}
