package repositories

import (
	"api/config"
	"api/models"
	"errors"

	"gorm.io/gorm"
)

type EmailQueueRepository struct {
}

func NewEmailQueueRepository() *EmailQueueRepository {
	return &EmailQueueRepository{}
}

func (r *EmailQueueRepository) Create(job *models.EmailQueueModel) error {
	if err := config.DB.Create(job).Error; err != nil {
		return err
	}

	return nil
}

func (r *EmailQueueRepository) Get(jobs *[]models.EmailQueueModel) error {
	if err := config.DB.Find(jobs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *EmailQueueRepository) Delete(ID int64) error {
	if err := config.DB.Where("id = ?", ID).Delete(&models.EmailQueueModel{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}
