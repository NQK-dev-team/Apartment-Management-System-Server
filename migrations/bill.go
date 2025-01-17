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
	config.DB.AutoMigrate(&models.BillModel{})
	config.DB.AutoMigrate(&models.ExtraPaymentModel{})
	config.DB.Exec("ALTER TABLE bill ADD CONSTRAINT bill_status CHECK (status >= 1 AND status <= 4);")
	config.DB.Exec("ALTER TABLE bill ADD CONSTRAINT bill_payment CHECK ((status=2 AND payment_time IS NOT NULL AND payer_id IS NOT NULL) OR (NOT status=2 AND payment_time NOT NULL AND payer_id NOT NULL));")
}

func (m *BillMigration) Down() {
	config.DB.Migrator().DropTable(&models.ExtraPaymentModel{})
	config.DB.Migrator().DropTable(&models.BillModel{})
}
