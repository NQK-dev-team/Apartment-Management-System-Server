package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
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

func (u *ManagerScheduleModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *ManagerScheduleModel) BeforeUpdate(tx *gorm.DB) error {
	// if tx.Statement.Changed("UpdatedAt", "UpdatedBy") {
	// 	return errors.New(config.GetMessageCode("CONCURRENCY_ERROR"))
	// }

	isQuiet, _ := tx.Get("isQuiet")
	if isQuiet != nil && isQuiet.(bool) {
		return nil
	}

	userID, _ := tx.Get("userID")
	if userID != nil {
		tx.Statement.SetColumn("updated_by", userID.(int64))
	}
	tx.Statement.SetColumn("updated_at", time.Now())

	return nil
}
