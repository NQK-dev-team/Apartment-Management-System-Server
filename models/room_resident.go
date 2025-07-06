package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type RoomResidentModel struct {
	DefaultModel
	FirstName               string         `json:"firstName" gorm:"column:first_name;type:varchar(255);not null;"`
	MiddleName              sql.NullString `json:"middleName" gorm:"column:middle_name;type:varchar(255);"`
	LastName                string         `json:"lastName" gorm:"column:last_name;type:varchar(255);not null;"`
	SSN                     sql.NullString `json:"ssn" gorm:"column:ssn;type:varchar(12);"`
	OldSSN                  sql.NullString `json:"oldSSN" gorm:"column:old_ssn;type:varchar(9);"`
	DOB                     string         `json:"dob" gorm:"column:dob;type:date;not null;"`
	POB                     sql.NullString `json:"pob" gorm:"column:pob;type:varchar(255);"`
	Email                   sql.NullString `json:"email" gorm:"column:email;type:varchar(255);"`
	Phone                   sql.NullString `json:"phone" gorm:"column:phone;type:varchar(10);"`
	Gender                  int            `json:"gender" gorm:"column:gender;type:int;not null;"`                                     // 1: Male, 2: Female, 3: Other
	RelationWithHouseholder int            `json:"relationWithHouseholder" gorm:"column:relation_with_householder;type:int;not null;"` // 1: Child, 2: Spouse, 3: Parent, 4: Other
	UserAccountID           sql.NullInt64  `json:"userAccountID" gorm:"column:user_account_id;"`
	UserAccount             UserModel      `json:"userAccount" gorm:"foreignKey:user_account_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *RoomResidentModel) TableName() string {
	return "room_resident"
}

func (u *RoomResidentModel) BeforeCreate(tx *gorm.DB) error {
	userID, _ := tx.Get("userID")
	if userID != nil {
		u.CreatedBy = userID.(int64)
		u.UpdatedBy = userID.(int64)
	}
	// u.CreatedAt = time.Now()
	// u.UpdatedAt = time.Now()

	return nil
}

func (u *RoomResidentModel) BeforeUpdate(tx *gorm.DB) error {
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
