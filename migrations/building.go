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
	config.DB.AutoMigrate(&models.BuildingModel{})
	config.DB.AutoMigrate(&models.BuildingServiceModel{})
	config.DB.Exec("ALTER TABLE building_service ADD CONSTRAINT building_service_composite_key UNIQUE (id, building_id);")
	config.DB.AutoMigrate(&models.BuildingImageModel{})
	config.DB.Exec("ALTER TABLE building_image ADD CONSTRAINT building_image_composite_key UNIQUE (id, building_id);")
	config.DB.AutoMigrate(&models.ManagerScheduleModel{})
	config.DB.Exec("ALTER TABLE manager_schedule ADD CONSTRAINT manager_schedule_period CHECK (start_date<=end_date);")
	m.roomMigration.Up()
}

func (m *BuildingMigration) Down() {
	m.roomMigration.Down()
	config.DB.Migrator().DropTable(&models.ManagerScheduleModel{})
	config.DB.Migrator().DropTable(&models.BuildingImageModel{})
	config.DB.Migrator().DropTable(&models.BuildingServiceModel{})
	config.DB.Migrator().DropTable(&models.BuildingModel{})
}
