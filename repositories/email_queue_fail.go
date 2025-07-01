package repositories

import (
	"api/config"
	"api/models"
)

type EmailQueueFailRepository struct {
}

func NewEmailQueueFailRepository() *EmailQueueFailRepository {
	return &EmailQueueFailRepository{}
}

func (r *EmailQueueFailRepository) Create(failedJob *models.EmailQueueFailModel) error {
	if err := config.DBNoLog.Create(failedJob).Error; err != nil {
		return err
	}

	return nil
}
