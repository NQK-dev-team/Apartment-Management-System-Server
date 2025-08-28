package structs

type BaseWSStruct struct {
	Type int `json:"type"`
}

type NotificationWS struct {
	BaseWSStruct
	Users []int64 `json:"users"`
}
