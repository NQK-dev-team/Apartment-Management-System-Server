package models

type NotificationReceiverModel struct {
	NotifcationID   int64             `json:"notificationID" gorm:"primaryKey;column:notification_id;"`
	Notification    NotificationModel `json:"notification" gorm:"foreignKey:notification_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID          int64             `json:"userID" gorm:"primaryKey;column:user_id;"`
	User            UserModel         `json:"user" gorm:"foreignKey:user_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MarkAsRead      int               `json:"markAsRead" gorm:"column:mark_as_read;not null;default:0;"`
	MarkAsImportant int               `json:"markAsImportant" gorm:"column:mark_as_important;not null;default:0;"`
}

func (u *NotificationReceiverModel) TableName() string {
	return "notification_receiver"
}
