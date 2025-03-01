package models

type BuildingServiceModel struct {
	DefaultModel
	BuildingID int64 `json:"buildingID" gorm:"column:building_id;not null;"`
	// Building   BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name  string  `json:"name" gorm:"column:name;type:varchar(255);not null;"`
	Price float64 `json:"price" gorm:"column:price;type:numeric;not null;"`
}

func (u *BuildingServiceModel) TableName() string {
	return "building_service"
}
