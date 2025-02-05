package structs

type BuildingID struct {
	ID string `json:"id" validate:"required"`
}
