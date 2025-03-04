package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuildingRepository struct {
}

func NewBuildingRepository() *BuildingRepository {
	return &BuildingRepository{}
}

func (r *BuildingRepository) Get(ctx *gin.Context, building *[]models.BuildingModel) error {
	if err := config.DB.Model(&models.BuildingModel{}).Preload("Images").Find(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetBuildingBaseOnSchedule(ctx *gin.Context, building *[]models.BuildingModel, userID int64) error {
	if err := config.DB.Model(&models.BuildingModel{}).Preload("Images").Select("building.*").Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id").Where("manager_schedule.start_date <= now() AND manager_schedule.end_date >= now() AND manager_schedule.manager_id = ?", userID).Scan(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetById(ctx *gin.Context, building *models.BuildingModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("Images").Preload("Services").Preload("Rooms").First(building).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *BuildingRepository) GetNewID(ctx *gin.Context) (int64, error) {
	lastestBuilding := models.BuildingModel{}
	if err := config.DB.Order("id desc").Unscoped().First(&lastestBuilding).Error; err != nil {
		return 0, err
	}
	return lastestBuilding.ID + 1, nil
}

func (r *BuildingRepository) GetNewImageID(ctx *gin.Context) (int64, error) {
	lastestImage := models.BuildingImageModel{}
	if err := config.DB.Order("id desc").Unscoped().First(&lastestImage).Error; err != nil {
		return 0, err
	}
	return lastestImage.ID + 1, nil
}

func (r *BuildingRepository) GetNewServiceID(ctx *gin.Context) (int64, error) {
	lastestService := models.BuildingServiceModel{}
	if err := config.DB.Order("id desc").Unscoped().First(&lastestService).Error; err != nil {
		return 0, err
	}
	return lastestService.ID + 1, nil
}

func (r *BuildingRepository) Create(ctx *gin.Context, building *models.BuildingModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Create(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) Update(ctx *gin.Context, building *models.BuildingModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Save(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) Delete(ctx *gin.Context, building *models.BuildingModel) error {
	if err := config.DB.Delete(building).Error; err != nil {
		return err
	}
	return nil
}
