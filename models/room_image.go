package models

type RoomImageModel struct {
	DefaultFileModel
	RoomID     int64     `json:"roomID" gorm:"column:room_id;primaryKey;"`
	BuildingID int64     `json:"buildingID" gorm:"column:building_id;primaryKey;"`
	// Room       RoomModel `json:"room" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *RoomImageModel) TableName() string {
	return "room_image"
}
