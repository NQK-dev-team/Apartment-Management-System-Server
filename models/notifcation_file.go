package models

type NotificationFileModel struct {
	DefaultFileModel
	NotificationID int64             `json:"notificationID" gorm:"column:notification_id;primaryKey;"`
	Notification   NotificationModel `json:"notification" gorm:"foreignKey:notification_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *NotificationFileModel) TableName() string {
	return "notification_file"
}
