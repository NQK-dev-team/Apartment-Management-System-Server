package structs

type BuildingID struct {
	ID int64 `json:"id" validate:"required"`
}
