package models

import (
	"database/sql"
	"time"
)

type ManagerScheduleModel struct {
	DefaultModel
	StartDate  time.Time     `json:"startDate" gorm:"column:start_date;type:date;not null;default:now();"`
	EndDate    sql.NullTime  `json:"endDate" gorm:"column:end_date;type:date;"`
	ManagerID  int64         `json:"managerID" gorm:"column:manager_id;not null;"`
	Manager    UserModel     `json:"manager" gorm:"foreignKey:manager_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BuildingID int64         `json:"buildingID" gorm:"column:building_id;not null;"`
	Building   BuildingModel `json:"building" gorm:"foreignKey:building_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ManagerScheduleModel) TableName() string {
	return "manager_schedule"
}
