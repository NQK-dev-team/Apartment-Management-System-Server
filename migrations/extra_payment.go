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
	config.DB.Exec("ALTER TABLE extra_payment ADD CONSTRAINT extra_payment_composite_key UNIQUE (bill_id, contract_id);")
}

func (m *ExtraPaymentMigration) Down() {
	config.DB.Migrator().DropTable(&models.ExtraPaymentModel{})
}
