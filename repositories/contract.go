package repositories

import (
	"api/config"
	"api/models"
	"errors"
	"time"

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

func (r *ContractRepository) GetContractByIDs(ctx *gin.Context, contract *[]models.ContractModel, id []int64) error {
	if err := config.DB.Where("id in ?", id).Preload("Bills").Preload("SupportTickets").Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetContractByRoomID(ctx *gin.Context, contract *[]models.ContractModel, roomIDs []int64) error {
	if err := config.DB.Where("room_id in ?", roomIDs).Find(contract).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) Delete(ctx *gin.Context, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := config.DB.Set("isQuiet", true).Model(&models.ContractModel{}).Where("id in ?", id).UpdateColumns(models.ContractModel{
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

	if err := config.DB.Set("isQuiet", true).Model(&models.ContractFileModel{}).Where("contract_id in ?", id).UpdateColumns(models.ContractFileModel{
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
