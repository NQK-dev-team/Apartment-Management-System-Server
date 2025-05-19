package structs

import (
	"api/models"
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
	StartDate string `form:"startDate" validate:"required,datetime=2006-01-02"`
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
	StartDate string `form:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate   string `form:"endDate" validate:"required"`
}

type EditRoom struct {
	ID          int64   `form:"id" validate:"required"`
	No          int     `form:"no" validate:"required,min=1001"`
	Floor       int     `form:"floor" validate:"required,min=1"`
	Status      int     `form:"status" validate:"required,min=1,max=5"`
	Area        float64 `form:"area" validate:"required,gt=0"`
	Description string  `form:"description" validate:"required"`
	NewImages   []*multipart.FileHeader
	TotalImage  int `validate:"required,min=1"`
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
	DeletedRoomImages     []int64        `form:"deletedRoomImages[]"`
	NewRooms              []NewRoom      `form:"newRooms[]"`
	Rooms                 []EditRoom     `form:"rooms[]"`
	TotalFloor            int            `form:"totalFloor"`
	TotalImage            int            `validate:"required,min=1"`
}

type BuildingRoom struct {
	models.DefaultModel
	No              int                     `json:"no" gorm:"column:no;type:int;not null;"`
	Floor           int                     `json:"floor" gorm:"column:floor;type:int;not null;"`
	Description     string                  `json:"description" gorm:"column:description;type:varchar(255);"`
	Area            float64                 `json:"area" gorm:"column:area;type:numeric;not null;"`
	Status          int                     `json:"status" gorm:"column:status;type:int;not null;default:1;"` // 1: Rented, 2: Bought, 3: Available, 4: Maintenanced, 5: Unavailable
	BuildingID      int64                   `json:"buildingID" gorm:"column:building_id;not null;"`
	Contracts       []Contract              `json:"contracts" gorm:"foreignKey:room_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Images          []models.RoomImageModel `json:"images" gorm:"foreignKey:room_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BuildingName    string                  `json:"buildingName" gorm:"column:building_name;type:varchar(255);"`
	BuildingAddress string                  `json:"buildingAddress" gorm:"column:building_address;type:varchar(255);"`
}
