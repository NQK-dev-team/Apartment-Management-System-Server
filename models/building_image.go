package models

type BuildingImageModel struct {
	DefaultFileModel
	BuildingID int64 `json:"buildingID" gorm:"column:building_id;not null;"`
	// Building   BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *BuildingImageModel) TableName() string {
	return "building_image"
}
