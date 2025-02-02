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
	if err := config.DB.Find(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) GetById(ctx *gin.Context, building *models.BuildingModel, id int64) error {
	if err := config.DB.Where("id = ?", id).First(building).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
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
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Delete(building).Error; err != nil {
		return err
	}
	return nil
}
