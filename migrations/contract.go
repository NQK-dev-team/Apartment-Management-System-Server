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
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_status CHECK (status >= 1 AND status <= 5);")
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_sign_date CHECK ((status=4 AND sign_date IS NULL) OR (NOT status=4 AND sign_date IS NOT NULL));")
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_period CHECK (start_date<=end_date);")
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_type CHECK (type=1 OR type=2);")
	config.DB.Exec("ALTER TABLE contract ADD CONSTRAINT contract_buy CHECK ((type=2 AND end_date IS NULL) OR type=1);")
	config.DB.AutoMigrate(&models.ContractFileModel{})
	// config.DB.Exec("ALTER TABLE contract_file ADD CONSTRAINT contract_file_composite_key UNIQUE (id, contract_id);")
	m.billMigration.Up()
	config.DB.AutoMigrate(&models.RoomResidentModel{})
	config.DB.Exec("ALTER TABLE room_resident ADD CONSTRAINT room_resident_relationship CHECK ((relation_with_householder >= 1 AND relation_with_householder <= 3) OR relation_with_householder IS NULL);")
	config.DB.AutoMigrate(&models.RoomResidentListModel{})
	// config.DB.Exec("ALTER TABLE room_resident_list ADD CONSTRAINT room_resident_list_composite_key UNIQUE (contract_id, resident_id);")
}

func (m *ContractMigration) Down() {
	config.DB.Migrator().DropTable(&models.RoomResidentListModel{})
	config.DB.Migrator().DropTable(&models.RoomResidentModel{})
	m.billMigration.Down()
	config.DB.Migrator().DropTable(&models.ContractFileModel{})
	config.DB.Migrator().DropTable(&models.ContractModel{})
}
