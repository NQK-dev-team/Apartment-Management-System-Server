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
	config.DB.AutoMigrate(&models.BuildingImageModel{})
	config.DB.AutoMigrate(&models.ManagerScheduleModel{})
	m.roomMigration.Up()
}

func (m *BuildingMigration) Down() {
	m.roomMigration.Down()
	config.DB.Migrator().DropTable(&models.ManagerScheduleModel{})
	config.DB.Migrator().DropTable(&models.BuildingImageModel{})
	config.DB.Migrator().DropTable(&models.BuildingServiceModel{})
	config.DB.Migrator().DropTable(&models.BuildingModel{})
}
