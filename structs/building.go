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
	Description string                  `form:"description" validate:"required,max=255"`
	Images      []*multipart.FileHeader `validate:"required,min=1"`
}

type NewSchedule struct {
	ManagerID int64  `form:"managerID" validate:"required"`
	StartDate string `form:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate   string `form:"endDate" validate:"required"`
}

type NewService struct {
	Name  string  `form:"name" validate:"required,max=255"`
	Price float64 `form:"price" validate:"required,gt=0"`
}

type NewBuilding struct {
	Name       string                  `form:"name" validate:"required,max=255"`
	Address    string                  `form:"address" validate:"required,max=255"`
	TotalRoom  int                     `form:"totalRoom"`
	TotalFloor int                     `form:"totalFloor"`
	Services   []NewService            `form:"services[]"`
	Images     []*multipart.FileHeader `validate:"required,min=1"`
	Rooms      []NewRoom               `form:"rooms[]"`
	Schedules  []NewSchedule           `form:"schedules[]"`
}

type EditService struct {
	ID    int64   `form:"id" validate:"required"`
	Name  string  `form:"name" validate:"required,max=255"`
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
	Description string  `form:"description" validate:"required,max=255"`
	NewImages   []*multipart.FileHeader
	TotalImage  int `validate:"required,min=1"`
}

type EditRoom2 struct {
	Status            int     `form:"status" validate:"required,min=1,max=5"`
	Area              float64 `form:"area" validate:"required,gt=0"`
	Description       string  `form:"description" validate:"required,max=255"`
	NewRoomImages     []*multipart.FileHeader
	DeletedRoomImages []int64 `form:"deletedRoomImages[]"`
	TotalImage        int     `validate:"required,min=1"`
}

type EditBuilding struct {
	ID                    int64   `form:"id" validate:"required"`
	Name                  string  `form:"name" validate:"required,max=255"`
	Address               string  `form:"address" validate:"required,max=255"`
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

type AllBuildingStatistic struct {
	TotalBuildings         int64 `json:"totalBuildings" gorm:"column:total_buildings"`
	TotalRooms             int64 `json:"totalRooms" gorm:"column:total_rooms"`
	TotalRentedRooms       int64 `json:"totalRentedRooms" gorm:"column:total_rented_rooms"`
	TotalBoughtRooms       int64 `json:"totalBoughtRooms" gorm:"column:total_bought_rooms"`
	TotalAvailableRooms    int64 `json:"totalAvailableRooms" gorm:"column:total_available_rooms"`
	TotalMaintenancedRooms int64 `json:"totalMaintenancedRooms" gorm:"column:total_maintenanced_rooms"`
	TotalUnavailableRooms  int64 `json:"totalUnavailableRooms" gorm:"column:total_unavailable_rooms"`
}

type BuildingStatistic struct {
	TotalRooms                int64                    `json:"totalRooms" gorm:"column:total_rooms"`
	TotalRentedRooms          int64                    `json:"totalRentedRooms" gorm:"column:total_rented_rooms"`
	TotalBoughtRooms          int64                    `json:"totalBoughtRooms" gorm:"column:total_bought_rooms"`
	TotalAvailableRooms       int64                    `json:"totalAvailableRooms" gorm:"column:total_available_rooms"`
	TotalMaintenancedRooms    int64                    `json:"totalMaintenancedRooms" gorm:"column:total_maintenanced_rooms"`
	TotalUnavailableRooms     int64                    `json:"totalUnavailableRooms" gorm:"column:total_unavailable_rooms"`
	TotalContract             int64                    `json:"total" gorm:"column:total_contract"`
	TotalRent                 int64                    `json:"total_rent" gorm:"column:total_rent_contract"`
	TotalBuy                  int64                    `json:"total_buy" gorm:"column:total_buy_contract"`
	TotalActiveRent           int64                    `json:"total_active_rent" gorm:"column:total_active_rent_contract"`
	TotalActiveBuy            int64                    `json:"total_active_buy" gorm:"column:total_active_buy_contract"`
	TotalExpireRent           int64                    `json:"total_expire_rent" gorm:"column:total_expire_rent_contract"`
	TotalCancelRent           int64                    `json:"total_cancel_rent" gorm:"column:total_cancel_rent_contract"`
	TotalCancelBuy            int64                    `json:"total_cancel_buy" gorm:"column:total_cancel_buy_contract"`
	TotalWaitForSignatureRent int64                    `json:"total_wait_for_signature_rent" gorm:"column:total_wait_for_signature_rent_contract"`
	TotalWaitForSignatureBuy  int64                    `json:"total_wait_for_signature_buy" gorm:"column:total_wait_for_signature_buy_contract"`
	TotalNotInEffectRent      int64                    `json:"total_not_in_effect_rent" gorm:"column:total_not_in_effect_rent_contract"`
	TotalNotInEffectBuy       int64                    `json:"total_not_in_effect_buy" gorm:"column:total_not_in_effect_buy_contract"`
	TotalBill                 int                      `json:"totalBill" gorm:"column:total_bill"`
	TotalPaid                 int                      `json:"totalPaid" gorm:"column:total_paid"`
	TotalUnpaid               int                      `json:"totalUnpaid" gorm:"column:total_unpaid"`
	TotalOverdue              int                      `json:"totalOverdue" gorm:"column:total_overdue"`
	RevenueStatistic          []RevenueStatisticStruct `json:"revenueStatistic" gorm:"-"`
}
