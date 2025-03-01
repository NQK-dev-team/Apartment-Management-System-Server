package models

type ContractFileModel struct {
	DefaultFileModel
	ContractID int64         `json:"contract_id" gorm:"column:contract_id;not null;"`
	// Contract   ContractModel `json:"contract" gorm:"foreignKey:contract_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *ContractFileModel) TableName() string {
	return "contract_file"
}
