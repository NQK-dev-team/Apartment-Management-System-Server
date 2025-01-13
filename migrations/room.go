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
	config.DB.AutoMigrate(&models.RoomImageModel{})
	config.DB.Exec("ALTER TABLE room ADD CONSTRAINT room_status CHECK (status >= 0 AND status <= 4);")
	// config.DB.Exec("ALTER TABLE room ADD CONSTRAINT unique_room_id_building_id UNIQUE (ID, building_id);")
}

func (m *RoomMigration) Down() {
	config.DB.Migrator().DropTable(&models.RoomImageModel{})
	config.DB.Migrator().DropTable(&models.RoomModel{})
}
