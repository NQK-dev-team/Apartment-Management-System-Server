package migrations

import (
	"api/config"
	"api/models"
)

type BillPaymentMigration struct {
}

func NewBillPaymentMigration() *BillPaymentMigration {
	return &BillPaymentMigration{}
}

func (m *BillPaymentMigration) Up() {
	config.MigrationDB.AutoMigrate(&models.BillPaymentModel{})
	// config.DB.Exec("ALTER TABLE extra_payment ADD CONSTRAINT extra_payment_composite_key UNIQUE (bill_id, contract_id);")
}

func (m *BillPaymentMigration) Down() {
	config.MigrationDB.Migrator().DropTable(&models.BillPaymentModel{})
}
