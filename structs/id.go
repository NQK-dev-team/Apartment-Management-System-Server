package structs

type IDList struct {
	IDs []int64 `json:"IDs" validate:"required"`
}
