package structs

import "mime/multipart"

type BuildingID struct {
	ID int64 `json:"id" validate:"required"`
}

type NewManagerSchedule struct {
	ID        int64  `form:"ID" validate:"required"`
	StartDate string `form:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate   string `form:"endDate" validate:"datetime=2006-01-02,gtefield=StartDate"`
}

type NewRoom struct {
	No          int                     `form:"no" validate:"required,min=1001"`
	Floor       int                     `form:"floor" validate:"required,min=1"`
	Status      int                     `form:"status" validate:"require,min=1,max=5"`
	Area        float64                 `form:"area" validate:"required,gt=0"`
	Description string                  `form:"description" validate:"required"`
	Images      *[]multipart.FileHeader `validate:"required,min=1"`
}

type NewService struct {
	Name  string  `form:"name" validate:"required"`
	Price float64 `form:"price" validate:"required,gt=0"`
}

type NewBuilding struct {
	Name       string                  `form:"name" validate:"required"`
	Address    string                  `form:"address" validate:"required"`
	TotalRoom  int                     `form:"totalRoom"`
	TotalFloor int                     `form:"totalFloor"`
	Services   []NewService            `form:"services[]" validate:"required"`
	Managers   []NewManagerSchedule    `form:"managers[]" validate:"required"`
	Images     *[]multipart.FileHeader `validate:"required,min=1"`
	Rooms      []NewRoom               `form:"rooms[]" validate:"required"`
}
