package migrations

import (
	"api/config"
	"api/models"
)

type BillMigration struct {
}

func NewBillMigration() *BillMigration {
	return &BillMigration{}
}

func (m *BillMigration) Up() {
	config.MigrationDB.AutoMigrate(&models.BillModel{})
	config.DB.Exec("ALTER TABLE bill ADD CONSTRAINT bill_status CHECK (status >= 1 AND status <= 5);")
	config.DB.Exec("ALTER TABLE bill ADD CONSTRAINT bill_payment CHECK ((status=2 AND payment_time IS NOT NULL AND payer_id IS NOT NULL) OR (NOT status=2 AND payment_time NOT NULL AND payer_id NOT NULL));")
	// config.DB.Exec("ALTER TABLE bill ADD CONSTRAINT bill_composite_key UNIQUE (id, contract_id);")
	config.MigrationDB.AutoMigrate(&models.BillPaymentModel{})
}

func (m *BillMigration) Down() {
	config.MigrationDB.Migrator().DropTable(&models.BillPaymentModel{})
	config.MigrationDB.Migrator().DropTable(&models.BillModel{})
}
