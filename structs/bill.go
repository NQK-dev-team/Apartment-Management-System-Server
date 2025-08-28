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
	ID     int64   `json:"ID" validate:"required"`
	Amount float64 `json:"amount" validate:"required,min=0"`
	Name   string  `json:"name" validate:"required,max=255"`
	Note   string  `json:"note" validate:"omitempty,max=255"`
}

type NewPayment struct {
	Amount float64 `json:"amount" validate:"required,min=0"`
	Name   string  `json:"name" validate:"required,max=255"`
	Note   string  `json:"note" validate:"omitempty,max=255"`
}

type UpdateBill struct {
	Title           string       `json:"title" validate:"required,max=255"`
	Period          string       `json:"period" validate:"required,datetime=2006-01"`
	Status          int          `json:"status" validate:"required,min=1,max=5"`
	Note            string       `json:"note" validate:"omitempty,max=255"`
	Payments        []OldPayment `json:"payments" validate:"dive"`
	NewPayments     []NewPayment `json:"newPayments" validate:"dive"`
	DeletedPayments []int64      `json:"deletedPayments"`
	PayerID         int64        `json:"payerID" validate:"required_if=Status 2"`
	PaymentTime     string       `json:"paymentTime" validate:"required_unless=PayerID 0,omitempty,datetime=2006-01-02,validate_payment_time"`
}

type AddBill struct {
	Title        string       `json:"title" validate:"required,max=255"`
	Period       string       `json:"period" validate:"required,datetime=2006-01"`
	Status       int          `json:"status" validate:"required,min=1,max=5"`
	Note         string       `json:"note" validate:"omitempty,max=255"`
	ContractID   int64        `json:"contractID" validate:"required"`
	PayerID      int64        `json:"payerID" validate:"required_if=Status 2"`
	PaymentTime  string       `json:"paymentTime" validate:"required_unless=PayerID 0,omitempty,datetime=2006-01-02,validate_payment_time"`
	BillPayments []NewPayment `json:"billPayments" validate:"min=1,dive"`
}

type UploadBill struct {
	Title       string `json:"title" validate:"required,max=255"`
	Period      string `json:"period" validate:"required,datetime=2006-01"`
	Status      int    `json:"status" validate:"required,min=1,max=5"`
	Note        string `json:"note" validate:"omitempty,max=255"`
	ContractID  int64  `json:"contractID" validate:"required"`
	PayerID     int64  `json:"payerID" validate:"required_if=Status 2"`
	PaymentTime string `json:"paymentTime" validate:"required_unless=PayerID 0,omitempty,datetime=2006-01-02,validate_payment_time"`
}
