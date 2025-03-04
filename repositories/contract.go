package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractRepository struct {
}

func NewContractRepository() *ContractRepository {
	return &ContractRepository{}
}

func (r *ContractRepository) GetById(ctx *gin.Context, contract *models.ContractModel, id int64) error {
	if err := config.DB.Where("id = ?", id).Preload("Files").Preload("Bills").Preload("SupportTickets").First(contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *ContractRepository) QuietUpdate(ctx *gin.Context, contract *models.ContractModel) error {
	if err := config.DB.Set("isQuiet", true).Session(&gorm.Session{FullSaveAssociations: true}).Save(contract).Error; err != nil {
		return err
	}
	return nil
}
