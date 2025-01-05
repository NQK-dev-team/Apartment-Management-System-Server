package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
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
		return err
	}
	return nil
}

func (r *BuildingRepository) Create(ctx *gin.Context, building *models.BuildingModel) error {
	if err := config.DB.Create(building).Error; err != nil {
		return err
	}
	return nil
}

func (r *BuildingRepository) Update(ctx *gin.Context, building *models.BuildingModel) error {
	if err := config.DB.Save(building).Error; err != nil {
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
