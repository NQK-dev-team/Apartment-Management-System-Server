package models

type MessageFileModel struct {
	DefaultFileModel
	MessageID int64 `json:"messageID" gorm:"column:message_id;not null;"`
	// Message   MessageModel `json:"message" gorm:"foreignKey:message_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *MessageFileModel) TableName() string {
	return "message_file"
}
