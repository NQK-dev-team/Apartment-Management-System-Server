package models

import (
	"database/sql"
)

type UserModel struct {
	DefaultModel
	FirstName        string       `json:"firstName" gorm:"column:first_name;type:varchar(255);not null;"`
	MiddleName       string       `json:"middleName" gorm:"column:middle_name;type:varchar(255);"`
	LastName         string       `json:"lastName" gorm:"column:last_name;type:varchar(255);not null;"`
	SSN              string       `json:"ssn" gorm:"column:ssn;type:varchar(12);not null;uniqueIndex:idx_ssn;"`
	OldSSN           string       `json:"oldSSN" gorm:"column:old_ssn;type:varchar(9);unique;"`
	DOB              string       `json:"dob" gorm:"column:dob;type:date;not null;"`
	POB              string       `json:"pob" gorm:"column:pob;type:varchar(255);"`
	Email            string       `json:"email" gorm:"column:email;type:varchar(255);not null;uniqueIndex:idx_email;"`
	Password         string       `json:"password" gorm:"column:password;type:varchar(255);not null;"`
	Phone            string       `json:"phone" gorm:"column:phone;type:varchar(10);not null;"`
	SSNFrontFilePath string       `json:"ssnFrontFilePath" gorm:"column:ssn_front_file_path;type:varchar(255);not null;"`
	SSNBackFilePath  string       `json:"ssnBackFilePath" gorm:"column:ssn_back_file_path;type:varchar(255);not null;"`
	ProfileFilePath  string       `json:"profileFilePath" gorm:"column:profile_file_path;type:varchar(255);not null;"`
	EmailVerifiedAt  sql.NullTime `json:"emailVerifiedAt" gorm:"column:email_verified_at;type:timestamp with time zone;"`
	IsOwner          bool         `json:"isOwner" gorm:"column:is_owner;type:bool;not null;"`
	IsManager        bool         `json:"isManager" gorm:"column:is_manager;type:bool;not null;"`
	IsCustomer       bool         `json:"isCustomer" gorm:"column:is_customer;type:bool;not null;"`
}

func (u *UserModel) TableName() string {
	return "user"
}
