package models

type EmailQueueFailModel struct {
	ID            int64  `gorm:"primaryKey; column:id; autoIncrement; not null;"`
	Subject       string `gorm:"column:subject;type:varchar(255);not null;"`
	Body          string `gorm:"column:body;type:text;not null;"`
	ReceiverEmail string `gorm:"column:receiver_email;type:varchar(255);not null;"`
	Error         string `gorm:"column:error;type:text;"`
}

func (u *EmailQueueFailModel) TableName() string {
	return "email_queue_fail"
}