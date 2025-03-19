package repositories

import (
	"api/config"
	"api/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuildingRepository struct {
}

func NewBuildingRepository() *BuildingRepository {
	return &BuildingRepository{}
}

func (r *BuildingRepository) Get(ctx *gin.Context, building *[]models.BuildingModel) error {
	if err := config.DB.Model(&models.BuildingModel{}).Preload("Images").Find(building).Order("id asc").Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetBuildingBaseOnSchedule(ctx *gin.Context, building *[]models.BuildingModel, userID int64) error {
	if err := config.DB.Model(&models.BuildingModel{}).Preload("Images").Joins("JOIN manager_schedule ON manager_schedule.building_id = building.id").Where("manager_schedule.start_date <= now() AND manager_schedule.end_date >= now() AND manager_schedule.manager_id = ?", userID).Find(building).Order("id asc").Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetById(ctx *gin.Context, building *models.BuildingModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("Rooms").Preload("Rooms.Images").Preload("Images").Preload("Services").First(building).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// func (r *BuildingRepository) GetNewImageNo(ctx *gin.Context, buildingID int64) (int, error) {
// 	lastestImage := models.BuildingImageModel{}
// 	if err := config.DB.Where("building_id = ?", buildingID).Order("no desc").Unscoped().First(&lastestImage).Error; err != nil {
// 		return 0, err
// 	}
// 	return lastestImage.No + 1, nil
// }

func (r *BuildingRepository) Create(ctx *gin.Context, tx *gorm.DB, building *models.BuildingModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Model(&models.BuildingModel{}).Omit("ID").Create(building).Error; err != nil {
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

func (r *BuildingRepository) Delete(ctx *gin.Context, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := config.DB.Set("isQuiet", true).Model(&models.BuildingModel{}).Where("id in ?", id).UpdateColumns(models.BuildingModel{
		DefaultModel: models.DefaultModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}

	if err := config.DB.Set("isQuiet", true).Model(&models.BuildingImageModel{}).Where("building_id in ?", id).UpdateColumns(models.BuildingImageModel{
		DefaultFileModel: models.DefaultFileModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}

	if err := config.DB.Set("isQuiet", true).Model(&models.BuildingServiceModel{}).Where("building_id in ?", id).UpdateColumns(models.BuildingServiceModel{
		DefaultModel: models.DefaultModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) AddImage(ctx *gin.Context, tx *gorm.DB, image *[]models.BuildingImageModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BuildingImageModel{}).Omit("ID").Create(image).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) DeleteImages(ctx *gin.Context, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := config.DB.Set("isQuiet", true).Model(&models.BuildingImageModel{}).Where("building_id in ?", id).UpdateColumns(models.BuildingImageModel{
		DefaultFileModel: models.DefaultFileModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) AddServices(ctx *gin.Context, tx *gorm.DB, services *[]models.BuildingServiceModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.BuildingServiceModel{}).Omit("ID").Create(services).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) DeleteServices(ctx *gin.Context, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := config.DB.Set("isQuiet", true).Model(&models.BuildingServiceModel{}).Where("id in ?", id).UpdateColumns(models.BuildingServiceModel{
		DefaultModel: models.DefaultModel{
			DeletedBy: userID,
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
		},
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) EditService(ctx *gin.Context, service *models.BuildingServiceModel) error {
	userID := ctx.GetInt64("userID")
	if err := config.DB.Set("userID", userID).Save(service).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetServiceByID(ctx *gin.Context, service *models.BuildingServiceModel, id int64) error {
	if err := config.DB.Where("id = ?", id).First(service).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetBuildingSchedule(ctx *gin.Context, buildingID int64, schedule *[]models.ManagerScheduleModel) error {
	if err := config.DB.Preload("Manager").Where("building_id = ?", buildingID).Find(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetManagerBuildingSchedule(ctx *gin.Context, buildingID int64, schedule *[]models.ManagerScheduleModel, mangerID int64) error {
	if err := config.DB.Preload("Manager").Where("building_id = ? AND manager_id = ?", buildingID, mangerID).Find(schedule).Error; err != nil {
		return err
	}
	return nil
}
