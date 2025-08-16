package repositories

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/structs"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) GetReceiverNotificationRelation(ctx *gin.Context, notificationID int64, userID int64, relation *models.NotificationReceiverModel) error {
	return config.DB.Model(&models.NotificationReceiverModel{}).Where("notification_id = ? AND user_id = ?", notificationID, userID).Find(relation).Error
}

func (r *NotificationRepository) UpdateReceiverNotificationRelation(ctx *gin.Context, tx *gorm.DB, relation *models.NotificationReceiverModel) error {
	userID := ctx.GetInt64("userID")
	return tx.Set("userID", userID).Model(&models.NotificationReceiverModel{}).Where("notification_id = ? AND user_id = ?", relation.NotificationID, userID).Save(relation).Error
}

func (r *NotificationRepository) GetNotificationByID(ctx *gin.Context, id int64, notification *models.NotificationModel) error {
	return config.DB.Model(&models.NotificationModel{}).Preload("Files").Where("id = ?", id).Find(notification).Error
}

func (r *NotificationRepository) DeleteNotification(ctx *gin.Context, tx *gorm.DB, id int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.NotificationModel{}).Where("id = ?", id).UpdateColumns(models.NotificationModel{
		DefaultModel: models.DefaultModel{
			DeletedAt: gorm.DeletedAt{
				Valid: true,
				Time:  now,
			},
			DeletedBy: userID,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetSentNotification(ctx *gin.Context, notifications *[]models.NotificationModel, userID, limit, offset int64) error {
	return config.DB.Model(&models.NotificationModel{}).Preload("Files").Where("sender_id = ?", userID).Order("send_time DESC").Limit(int(limit)).Offset(int(offset)).Find(notifications).Error
}

func (r *NotificationRepository) GetInboxNotification(ctx *gin.Context, notifications *[]structs.Notification, userID, limit, offset int64) error {
	return config.DB.Model(&models.NotificationModel{}).Preload("Sender").Preload("Files").Distinct().Select("notification.*, notification_receiver.mark_as_read as is_read, notification_receiver.mark_as_important as is_marked, \"user\".first_name as sender_first_name, \"user\".last_name as sender_last_name, \"user\".middle_name as sender_middle_name").
		Joins("JOIN notification_receiver ON notification_receiver.notification_id = notification.id").
		Joins("JOIN \"user\" ON \"user\".id = notification.sender_id").
		Where("notification_receiver.user_id = ?", userID).Order("notification.send_time DESC").Limit(int(limit)).Offset(int(offset)).Find(notifications).Error
}

func (r *NotificationRepository) GetMarkedNotification(ctx *gin.Context, notifications *[]structs.Notification, userID, limit, offset int64) error {
	return config.DB.Model(&models.NotificationModel{}).Preload("Sender").Preload("Files").Distinct().Select("notification.*, notification_receiver.mark_as_read as is_read, notification_receiver.mark_as_important as is_marked, \"user\".first_name as sender_first_name, \"user\".last_name as sender_last_name, \"user\".middle_name as sender_middle_name").
		Joins("JOIN notification_receiver ON notification_receiver.notification_id = notification.id").
		Joins("JOIN \"user\" ON \"user\".id = notification.sender_id").
		Where("notification_receiver.user_id = ? AND is_marked = ?", userID, constants.Common.Notification.MarkedStatus).Order("notification.send_time DESC").Limit(int(limit)).Offset(int(offset)).Find(notifications).Error
}

func (r *NotificationRepository) CreateNotification(ctx *gin.Context, tx *gorm.DB, notification *models.NotificationModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.NotificationModel{}).Omit("ID").Create(notification).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) AddNotificationReceivers(ctx *gin.Context, tx *gorm.DB, receivers *[]models.NotificationReceiverModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.NotificationReceiverModel{}).Create(receivers).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) AddNotificationFiles(ctx *gin.Context, tx *gorm.DB, files *[]models.NotificationFileModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.NotificationFileModel{}).Omit("ID").Save(&files).Error; err != nil {
		return err
	}
	return nil
}
