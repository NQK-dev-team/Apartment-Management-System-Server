package models

import (
	"database/sql"
	"time"
)

type ManagerScheduleModel struct {
	DefaultModel
	StartDate  time.Time     `json:"start_date" gorm:"column:start_date;type:date;not null;default:now();"`
	EndDate    sql.NullTime  `json:"end_date" gorm:"column:end_date;type:date;"`
	ManagerID  int64         `json:"manager_id" gorm:"column:manager_id;not null;"`
	Manager    UserModel     `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BuildingID int64         `json:"building_id" gorm:"column:building_id;not null;"`
	Building   BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ManagerScheduleModel) TableName() string {
	return "manager_schedule"
}
