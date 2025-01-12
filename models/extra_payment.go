package models

type ExtraPaymentModel struct {
	DefaultModel
	Name       string    `json:"name" gorm:"column:name;not null;type:varchar(255);"`
	Amount     float64   `json:"amount" gorm:"column:amount;not null;type:numeric;"`
	Note       string    `json:"note" gorm:"column:note;type:varchar(255);"`
	ContractID int64     `json:"contractID" gorm:"column:contract_id;primaryKey;"`
	BillID     int64     `json:"billID" gorm:"column:bill_id;primaryKey;"`
	Bill       BillModel `json:"bill" gorm:"foreignKey:bill_id,contract_id;references:id,contract_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *ExtraPaymentModel) TableName() string {
	return "extra_payment"
}
