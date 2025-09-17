package structs

type JWTPayload struct {
	UserID       int64  `json:"userID"`
	FullName     string `json:"fullName"`
	ImagePath    string `json:"imagePath"`
	IsCustomer   bool   `json:"isCustomer"`
	IsManager    bool   `json:"isManager"`
	IsOwner      bool   `json:"isOwner"`
	UserNo       string `json:"userNo"`
	TicketByPass bool   `json:"ticketByPass"`
}

type JTWClaim struct {
	UserID       int64  `json:"userID"`
	FullName     string `json:"fullName"`
	ImagePath    string `json:"imagePath"`
	IsCustomer   bool   `json:"isCustomer"`
	IsManager    bool   `json:"isManager"`
	IsOwner      bool   `json:"isOwner"`
	UserNo       string `json:"userNo"`
	TicketByPass bool   `json:"ticketByPass"`
	ServiceToken string `json:"serviceToken"`
	IAT          int64  `json:"iat"`
	EXP          int64  `json:"exp"`
}
