package models

type BuildingModel struct {
	DefaultModel
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null;"`
	Address    string `json:"address" gorm:"column:address;type:varchar(255);not null;"`
	TotalFloor int    `json:"totalFloor" gorm:"column:total_floor;type:int;not null;"`
	TotalRoom  int    `json:"totalRoom" gorm:"column:total_room;type:int;not null;"`
}

func (u *BuildingModel) TableName() string {
	return "building"
}
