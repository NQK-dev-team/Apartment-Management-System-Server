package migrations

import (
	"api/config"
	"api/models"
)

type ContractMigration struct {
	billMigration *BillMigration
}

func NewContractMigration() *ContractMigration {
	return &ContractMigration{
		billMigration: NewBillMigration(),
	}
}

func (m *ContractMigration) Up() {
	config.DB.AutoMigrate(&models.ContractModel{})
	config.DB.AutoMigrate(&models.ContractFileModel{})
	m.billMigration.Up()
	config.DB.AutoMigrate(&models.RoomResidentModel{})
	config.DB.AutoMigrate(&models.RoomResidentListModel{})
}

func (m *ContractMigration) Down() {
	config.DB.Migrator().DropTable(&models.RoomResidentListModel{})
	config.DB.Migrator().DropTable(&models.RoomResidentModel{})
	m.billMigration.Down()
	config.DB.Migrator().DropTable(&models.ContractFileModel{})
	config.DB.Migrator().DropTable(&models.ContractModel{})
}
