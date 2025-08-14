package structs

type NotificationWS struct {
	// IsForAllStaffs    bool    `json:"isForAllStaffs"`
	// IsForAllCustomers bool    `json:"isForAllCustomers"`
	Users []int64 `json:"users"`
}
