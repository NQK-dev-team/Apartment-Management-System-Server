package migrations

import (
	"api/config"
	"api/models"
)

type BuildingMigration struct {
	roomMigration *RoomMigration
}

func NewBuildingMigration() *BuildingMigration {
	return &BuildingMigration{
		roomMigration: NewRoomMigration(),
	}
}

func (m *BuildingMigration) Up() {
	config.MigrationDB.AutoMigrate(&models.BuildingModel{})
	config.MigrationDB.AutoMigrate(&models.BuildingServiceModel{})
	// config.DB.Exec("ALTER TABLE building_service ADD CONSTRAINT building_service_composite_key UNIQUE (id, building_id);")
	config.MigrationDB.AutoMigrate(&models.BuildingImageModel{})
	// config.DB.Exec("ALTER TABLE building_image ADD CONSTRAINT building_image_composite_key UNIQUE (id, building_id);")
	config.MigrationDB.AutoMigrate(&models.ManagerScheduleModel{})
	config.DB.Exec("ALTER TABLE manager_schedule ADD CONSTRAINT manager_schedule_period CHECK (start_date<=end_date);")
	m.roomMigration.Up()
}

func (m *BuildingMigration) Down() {
	m.roomMigration.Down()
	config.MigrationDB.Migrator().DropTable(&models.ManagerScheduleModel{})
	config.MigrationDB.Migrator().DropTable(&models.BuildingImageModel{})
	config.MigrationDB.Migrator().DropTable(&models.BuildingServiceModel{})
	config.MigrationDB.Migrator().DropTable(&models.BuildingModel{})
}
