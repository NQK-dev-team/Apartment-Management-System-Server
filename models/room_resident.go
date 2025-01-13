package models

type RoomResidentModel struct {
	DefaultModel
	FirstName               string    `json:"firstName" gorm:"column:first_name;type:varchar(255);not null;"`
	MiddleName              string    `json:"middleName" gorm:"column:middle_name;type:varchar(255);"`
	LastName                string    `json:"lastName" gorm:"column:last_name;type:varchar(255);not null;"`
	SSN                     string    `json:"ssn" gorm:"column:ssn;type:varchar(12);not null;uniqueIndex:idx_ssn;"`
	OldSSN                  string    `json:"oldSSN" gorm:"column:old_ssn;type:varchar(9);unique;"`
	DOB                     string    `json:"dob" gorm:"column:dob;type:date;not null;"`
	POB                     string    `json:"pob" gorm:"column:pob;type:varchar(255);"`
	Email                   string    `json:"email" gorm:"column:email;type:varchar(255);not null;uniqueIndex:idx_email;"`
	Phone                   string    `json:"phone" gorm:"column:phone;type:varchar(10);not null;"`
	RelationWithHouseholder int       `json:"relationWithHouseholder" gorm:"column:relation_with_householder;type:int;"` // 0: Child, 1: Spouse, 2: Parent, null: Other
	UserAccountID           int64     `json:"userAccountID" gorm:"column:user_account_id;"`
	UserAccount             UserModel `json:"userAccount" gorm:"foreignKey:user_account_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *RoomResidentModel) TableName() string {
	return "room_resident"
}
