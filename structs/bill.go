package structs

import (
	"api/models"
	"database/sql"
	"time"
)

type Bill struct {
	models.DefaultModel
	Title                string                     `json:"title" gorm:"column:title;not null;type:varchar(255);"`
	Period               time.Time                  `json:"period" gorm:"column:period;not null;type:date;"`
	Status               int                        `json:"status" gorm:"column:status;not null;type:int;"` // 1: Unpaid, 2: Paid, 3: Overdue, 4: Processing
	Note                 string                     `json:"note" gorm:"column:note;type:varchar(255);"`
	PaymentTime          sql.NullTime               `json:"paymentTime" gorm:"column:payment_time;type:timestamp with time zone;default:now();"`
	Amount               float64                    `json:"amount" gorm:"column:amount;not null;type:numeric;"`
	TotalAmountWithExtra float64                    `json:"totalAmountWithExtra" gorm:"column:total_amount_with_extra;not null;type:numeric;"`
	PayerID              int64                      `json:"payerID" gorm:"column:payer_id;"`
	Payer                models.UserModel           `json:"payer" gorm:"foreignKey:payer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContractID           int64                      `json:"contractID" gorm:"column:contract_id;not null;"`
	ExtraPayments        []models.ExtraPaymentModel `json:"extraPayments" gorm:"foreignKey:bill_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Contract             models.ContractModel       `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	// ExtraPayments []ExtraPaymentModel `json:"extraPayments" gorm:"foreignKey:bill_id,contract_id;references:id,contract_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BuildingName    string `json:"buildingName" gorm:"building_name"`
	BuildingAddress string `json:"buildingAddress" gorm:"building_address"`
	RoomNo          int    `json:"roomNo" gorm:"room_no"`
	RoomFloor       int    `json:"roomFloor" gorm:"room_floor"`
}
