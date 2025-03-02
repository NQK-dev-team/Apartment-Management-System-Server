package structs

type BillID struct {
	ID int64 `json:"id" validate:"required"`
}

type NewBill struct {
}
