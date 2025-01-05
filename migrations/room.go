package migrations

import (
	"api/config"
	"api/models"
)

type RoomMigration struct {
}

func NewRoomMigration() *RoomMigration {
	return &RoomMigration{}
}

func (m *RoomMigration) Up() {
	config.DB.AutoMigrate(&models.RoomModel{})
}

func (m *RoomMigration) Down() {
	config.DB.Migrator().DropTable(&models.RoomModel{})
}
