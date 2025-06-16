package structs

import (
	"api/models"
	"database/sql"
	"time"
)

type Contract struct {
	models.DefaultModel
	Status        int                        `json:"status" gorm:"column:status;type:int;not null;"` // 1: Active, 2: Expired, 3: Cancelled, 4: Waiting for signatures, 5: Not in effect yet
	Value         float64                    `json:"value" gorm:"column:value;type:float;not null;"`
	Type          int                        `json:"type" gorm:"column:type;type:int;not null;"` // 1: Rent, 2: Buy
	StartDate     time.Time                  `json:"startDate" gorm:"column:start_date;type:date;not null;default:now();"`
	EndDate       sql.NullTime               `json:"endDate" gorm:"column:end_date;type:date;"`
	SignDate      sql.NullTime               `json:"signDate" gorm:"column:sign_date;type:date;default:now();"`
	CreatorID     int64                      `json:"creatorID" gorm:"column:creator_id;not null;"`
	Creator       models.UserModel           `json:"creator" gorm:"foreignKey:creator_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HouseholderID int64                      `json:"householderID" gorm:"column:householder_id;not null;"`
	Householder   models.UserModel           `json:"householder" gorm:"foreignKey:householder_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomID        int64                      `json:"roomID" gorm:"column:room_id;not null;"`
	Bills         []models.BillModel         `json:"bills" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Files         []models.ContractFileModel `json:"files" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Residents     []models.RoomResidentModel `json:"residents" gorm:"-"`
	// SupportTickets []SupportTicketModel `json:"supportTickets" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// BuildingID    int64     `json:"buildingID" gorm:"column:building_id;not null;"`
	// Room          RoomModel           `json:"room" gorm:"foreignKey:room_id,building_id;references:id,building_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BuildingName    string `json:"buildingName" gorm:"building_name"`
	BuildingAddress string `json:"buildingAddress" gorm:"building_address"`
	RoomNo          int    `json:"roomNo" gorm:"room_no"`
	RoomFloor       int    `json:"roomFloor" gorm:"room_floor"`
}
