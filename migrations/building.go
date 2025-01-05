package migrations

import (
	"api/config"
	"api/models"
)

type BuildingMigration struct {
}

func NewBuildingMigration() *BuildingMigration {
	return &BuildingMigration{}
}

func (m *BuildingMigration) Up() {
	config.DB.AutoMigrate(&models.BuildingModel{})
}

func (m *BuildingMigration) Down() {
	config.DB.Migrator().DropTable(&models.BuildingModel{})
}
