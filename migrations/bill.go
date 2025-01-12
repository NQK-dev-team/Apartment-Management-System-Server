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
}

func (m *BillMigration) Down() {
	config.DB.Migrator().DropTable(&models.ExtraPaymentModel{})
	config.DB.Migrator().DropTable(&models.BillModel{})
}
