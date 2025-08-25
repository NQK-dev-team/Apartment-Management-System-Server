package structs

import (
	"api/models"
	"database/sql"
	"time"
)

type Bill struct {
	models.DefaultModel
	Title           string                    `json:"title" gorm:"column:title;not null;type:varchar(255);"`
	Period          time.Time                 `json:"period" gorm:"column:period;not null;type:date;"`
	Status          int                       `json:"status" gorm:"column:status;not null;type:int;"` // 1: Unpaid, 2: Paid, 3: Overdue, 4: Processing
	Note            sql.NullString            `json:"note" gorm:"column:note;type:varchar(255);"`
	PaymentTime     sql.NullTime              `json:"paymentTime" gorm:"column:payment_time;type:timestamp with time zone;default:now();"`
	Amount          float64                   `json:"amount" gorm:"column:amount;not null;type:numeric;"`
	PayerID         sql.NullInt64             `json:"payerID" gorm:"column:payer_id;"`
	Payer           models.UserModel          `json:"payer" gorm:"foreignKey:payer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContractID      int64                     `json:"contractID" gorm:"column:contract_id;not null;"`
	BillPayments    []models.BillPaymentModel `json:"billPayments" gorm:"foreignKey:bill_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Contract        models.ContractModel      `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	BuildingName    string                    `json:"buildingName" gorm:"building_name"`
	BuildingAddress string                    `json:"buildingAddress" gorm:"building_address"`
	RoomNo          int                       `json:"roomNo" gorm:"room_no"`
	RoomFloor       int                       `json:"roomFloor" gorm:"room_floor"`
}

type OldPayment struct {
	ID     int64   `json:"ID" validation:"required"`
	Amount float64 `json:"amount" validation:"required,min=0"`
	Name   string  `json:"name" validation:"required"`
	Note   string  `json:"note"`
}

type NewPayment struct {
	Amount float64 `json:"amount" validation:"required,min=0"`
	Name   string  `json:"name" validation:"required"`
	Note   string  `json:"note"`
}

type UpdateBill struct {
	Title           string       `json:"title" validation:"required"`
	Period          string       `json:"period" validation:"required,datetime=2006-01"`
	Status          int          `json:"status" validation:"required,min=1,max=5"`
	Note            string       `json:"note"`
	Payments        []OldPayment `json:"payments" validation:"dive"`
	NewPayments     []NewPayment `json:"newPayments" validation:"dive"`
	DeletedPayments []int64      `json:"deletedPayments"`
	PayerID         int64        `json:"payerID" validation:"required_if=Status 2"`
	PaymentTime     string       `json:"paymentTime" validation:"required_unless=PayerID 0,datetime=2006-01-02,validate_payment_time"`
}

type AddBill struct {
	Title        string       `json:"title" validation:"required"`
	Period       string       `json:"period" validation:"required,datetime=2006-01"`
	Status       int          `json:"status" validation:"required,min=1,max=5"`
	Note         string       `json:"note"`
	ContractID   int64        `json:"contractID" validation:"required"`
	PayerID      int64        `json:"payerID" validation:"required_if=Status 2"`
	PaymentTime  string       `json:"paymentTime" validation:"required_unless=PayerID 0,datetime=2006-01-02,validate_payment_time"`
	BillPayments []NewPayment `json:"billPayments" validation:"min=1,dive"`
}

type UploadBill struct {
	Title       string `json:"title" validation:"required"`
	Period      string `json:"period" validation:"required,datetime=2006-01"`
	Status      int    `json:"status" validation:"required,min=1,max=5"`
	Note        string `json:"note"`
	ContractID  int64  `json:"contractID" validation:"required"`
	PayerID     int64  `json:"payerID" validation:"required_if=Status 2"`
	PaymentTime string `json:"paymentTime" validation:"required_unless=PayerID 0,datetime=2006-01-02"`
}
