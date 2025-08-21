package models

import (
	"database/sql"
	"time"
)

type UploadFileModel struct {
	ID        int64     `json:"ID" gorm:"primaryKey; column:id; autoIncrement; not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	CreatorID int64     `json:"creatorID" gorm:"column:creator_id;not null;"`
	Creator   UserModel `json:"creator" gorm:"foreignKey:creator_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FileName  string    `json:"fileName" gorm:"column:file_name;type:varchar(255);not null;"`
	URLPath   string    `json:"urlPath" gorm:"column:url_path;type:text;not null;"`
	// StoragePath   string        `json:"-" gorm:"column:storage_path;type:text;not null;"`
	Size          int64         `json:"size" gorm:"column:size;type:bigint;not null;"`
	UploadType    int           `json:"uploadType" gorm:"column:upload_type;type:int;not null;"` // 1: Customer, 2: Contract, 3: Bill
	ProcessDate   sql.NullTime  `json:"processDate" gorm:"column:process_date;type:date;"`
	ProcessResult sql.NullInt64 `json:"processResult" gorm:"column:process_result;type:int;"` // 1: Success, 2: Failed
}

func (u *UploadFileModel) TableName() string {
	return "upload_file"
}
