package structs

import (
	"mime/multipart"
)

type BuildingID struct {
	ID int64 `json:"id" validate:"required"`
}

type NewRoom struct {
	No          int                     `form:"no" validate:"required,min=1001"`
	Floor       int                     `form:"floor" validate:"required,min=1"`
	Status      int                     `form:"status" validate:"required,min=1,max=5"`
	Area        float64                 `form:"area" validate:"required,gt=0"`
	Description string                  `form:"description" validate:"required"`
	Images      []*multipart.FileHeader `validate:"required,min=1"`
}

type NewSchedule struct {
	ManagerID int64  `form:"managerID" validate:"required"`
	StartDate string `form:"startDate" validate:"required"`
	EndDate   string `form:"endDate" validate:"required"`
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
	Services   []NewService            `form:"services[]"`
	Images     []*multipart.FileHeader `validate:"required,min=1"`
	Rooms      []NewRoom               `form:"rooms[]"`
	Schedules  []NewSchedule           `form:"schedules[]"`
}

type EditService struct {
	ID    int64   `form:"id" validate:"required"`
	Name  string  `form:"name" validate:"required"`
	Price float64 `form:"price" validate:"required,gt=0"`
}

type EditSchedule struct {
	ID        int64  `form:"id" validate:"required"`
	ManagerID int64  `form:"managerID" validate:"required"`
	StartDate string `form:"startDate" validate:"required"`
	EndDate   string `form:"endDate" validate:"required"`
}

type EditRoom struct {
	ID          int64                   `form:"id" validate:"required"`
	No          int                     `form:"no" validate:"required,min=1001"`
	Floor       int                     `form:"floor" validate:"required,min=1"`
	Status      int                     `form:"status" validate:"required,min=1,max=5"`
	Area        float64                 `form:"area" validate:"required,gt=0"`
	Description string                  `form:"description" validate:"required"`
	Images      []*multipart.FileHeader `validate:"required,min=1"`
}

type EditBuilding struct {
	ID                    int64   `form:"id" validate:"required"`
	Name                  string  `form:"name" validate:"required"`
	Address               string  `form:"address" validate:"required"`
	DeletedBuildingImages []int64 `form:"deletedBuildingImages[]"`
	NewBuildingImages     []*multipart.FileHeader
	DeletedServices       []int64        `form:"deletedServices[]"`
	NewServices           []NewService   `form:"newServices[]"`
	Services              []EditService  `form:"services[]"`
	DeletedSchedules      []int64        `form:"deletedSchedules[]"`
	NewSchedules          []NewSchedule  `form:"newSchedules[]"`
	Schedules             []EditSchedule `form:"schedules[]"`
	DeletedRooms          []int64        `form:"deletedRooms[]"`
	NewRooms              []NewRoom      `form:"newRooms[]"`
	Rooms                 []EditRoom     `form:"rooms[]"`
	TotalFloor            int            `form:"totalFloor"`
}
