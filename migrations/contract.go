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
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_status CHECK (status >= 0 AND status <= 4);")
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_sign_date CHECK ((status=3 AND sign_date IS NULL) OR (NOT status=3 AND sign_date IS NOT NULL));")
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_buy CHECK (type=1 AND end_date IS NULL);")
	config.DB.Exec("ALTER TABLE room_resident ADD CONSTRAINT room_resident_relationship CHECK ((relation_with_householder >= 0 AND relation_with_householder <= 2) OR relation_with_householder IS NULL);")
}

func (m *ContractMigration) Down() {
	config.DB.Migrator().DropTable(&models.RoomResidentListModel{})
	config.DB.Migrator().DropTable(&models.RoomResidentModel{})
	m.billMigration.Down()
	config.DB.Migrator().DropTable(&models.ContractFileModel{})
	config.DB.Migrator().DropTable(&models.ContractModel{})
}
