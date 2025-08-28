package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationService struct {
	notificationRepository *repositories.NotificationRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepository: repositories.NewNotificationRepository(),
	}
}

func (s *NotificationService) UpdateNotificationReadStatus(ctx *gin.Context, id int64, isRead bool) (bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Owner {
		return false, nil
	}

	relation := &models.NotificationReceiverModel{}

	if err := s.notificationRepository.GetReceiverNotificationRelation(ctx, id, ctx.GetInt64("userID"), relation); err != nil {
		return true, err
	}

	if relation.NotificationID == 0 || relation.UserID == 0 {
		return false, nil
	}

	if isRead {
		relation.MarkAsRead = constants.Common.Notification.ReadStatus
	} else {
		relation.MarkAsRead = constants.Common.Notification.UnreadStatus
	}

	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.notificationRepository.UpdateReceiverNotificationRelation(ctx, tx, relation); err != nil {
			return err
		}

		return nil
	})
}

func (s *NotificationService) UpdateNotificationImportantStatus(ctx *gin.Context, id int64, isImportant bool) (bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Owner {
		return false, nil
	}

	relation := &models.NotificationReceiverModel{}

	if err := s.notificationRepository.GetReceiverNotificationRelation(ctx, id, ctx.GetInt64("userID"), relation); err != nil {
		return true, err
	}

	if relation.NotificationID == 0 || relation.UserID == 0 {
		return false, nil
	}

	if isImportant {
		relation.MarkAsImportant = constants.Common.Notification.MarkedStatus
	} else {
		relation.MarkAsImportant = constants.Common.Notification.UnmarkedStatus
	}

	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.notificationRepository.UpdateReceiverNotificationRelation(ctx, tx, relation); err != nil {
			return err
		}

		return nil
	})
}

func (s *NotificationService) MarkMultiNotiAsRead(ctx *gin.Context, ids *[]int64) (bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Owner {
		return false, nil
	}

	for _, id := range *ids {
		relation := &models.NotificationReceiverModel{}

		if err := s.notificationRepository.GetReceiverNotificationRelation(ctx, id, ctx.GetInt64("userID"), relation); err != nil {
			return true, err
		}

		if relation.NotificationID == 0 || relation.UserID == 0 {
			return false, nil
		}
	}

	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range *ids {
			relation := &models.NotificationReceiverModel{}

			if err := s.notificationRepository.GetReceiverNotificationRelation(ctx, id, ctx.GetInt64("userID"), relation); err != nil {
				return err
			}

			relation.MarkAsRead = constants.Common.Notification.ReadStatus

			if err := s.notificationRepository.UpdateReceiverNotificationRelation(ctx, tx, relation); err != nil {
				return err
			}
		}

		return nil
	})
}

// func (s *NotificationService) DeleteNotification(ctx *gin.Context, id int64, receivers *[]int64) (bool, error) {
// 	jwt := ctx.GetString("jwt")

// 	if !exists {
// 		return true, errors.New("jwt not found")
// 	}

// 	token, err := utils.ValidateJWTToken(jwt)

// 	if err != nil {
// 		return true, err
// 	}

// 	claim := &structs.JTWClaim{}

// 	utils.ExtractJWTClaim(token, claim)

// 	notification := &models.NotificationModel{}

// 	if err := s.notificationRepository.GetNotificationByID(ctx, id, notification); err != nil {
// 		return true, err
// 	}

// 	if notification.ID == 0 || notification.SenderID != claim.UserID {
// 		return false, nil
// 	}

// 	return true, config.DB.Transaction(func(tx *gorm.DB) error {
// 		if err := s.notificationRepository.DeleteNotification(ctx, tx, id, receivers); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

func (s *NotificationService) GetSentNotifications(ctx *gin.Context, notifications *[]models.NotificationModel, limit, offset int64) error {
	userID := ctx.GetInt64("userID")
	return s.notificationRepository.GetSentNotification(ctx, notifications, userID, limit, offset)
}

func (s *NotificationService) GetInboxNotifications(ctx *gin.Context, notifications *[]structs.Notification, limit, offset int64) error {
	userID := ctx.GetInt64("userID")
	return s.notificationRepository.GetInboxNotification(ctx, notifications, userID, limit, offset)
}

func (s *NotificationService) GetMarkedNotifications(ctx *gin.Context, notifications *[]structs.Notification, limit, offset int64) error {
	userID := ctx.GetInt64("userID")
	return s.notificationRepository.GetMarkedNotification(ctx, notifications, userID, limit, offset)
}

func (s *NotificationService) AddNotification(ctx *gin.Context, newNotification *structs.NewNotification) error {
	userID := ctx.GetInt64("userID")

	deleteFileList := []string{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		notification := &models.NotificationModel{
			Title:    newNotification.Title,
			Content:  newNotification.Content,
			SenderID: userID,
			SendTime: time.Now(),
		}

		if err := s.notificationRepository.CreateNotification(ctx, tx, notification); err != nil {
			return err
		}

		receivers := []models.NotificationReceiverModel{}

		for _, receiverID := range newNotification.Receivers {
			receivers = append(receivers, models.NotificationReceiverModel{
				NotificationID: notification.ID,
				UserID:         receiverID,
			})
		}

		if err := s.notificationRepository.AddNotificationReceivers(ctx, tx, &receivers); err != nil {
			return err
		}

		if len(newNotification.Files) > 0 {
			notificationIDStr := strconv.Itoa(int(notification.ID))

			files := []models.NotificationFileModel{}
			for _, file := range newNotification.Files {
				filePath, err := utils.StoreFile(file, constants.GetNotificationFileURL("files", notificationIDStr, ""))
				if err != nil {
					return err
				}
				files = append(files, models.NotificationFileModel{
					NotificationID: notification.ID,
					DefaultFileModel: models.DefaultFileModel{
						Path:  filePath,
						Title: filepath.Base(filePath),
					},
				})
				deleteFileList = append(deleteFileList, filePath)
			}

			if err := s.notificationRepository.AddNotificationFiles(ctx, tx, &files); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		for _, path := range deleteFileList {
			utils.RemoveFile(path)
		}
		return err
	}

	return nil
}

func (s *NotificationService) CheckUserGetNotification(ctx *gin.Context, notificationID int64) (bool, error) {
	userID := ctx.GetInt64("userID")

	role := ctx.GetString("role")

	relation := &models.NotificationReceiverModel{}

	if err := s.notificationRepository.GetReceiverNotificationRelation(ctx, notificationID, userID, relation); err != nil {
		return false, err
	}

	if role == constants.Roles.Customer {
		return relation.NotificationID != 0 && relation.UserID != 0, nil
	}

	isValid := false

	if role == constants.Roles.Manager && relation.NotificationID != 0 && relation.UserID != 0 {
		isValid = true
	}

	if isValid {
		return true, nil
	}

	notification := &models.NotificationModel{}

	if err := s.notificationRepository.GetNotificationByID(ctx, notificationID, notification); err != nil {
		return false, err
	}

	isValid = notification.SenderID == userID

	return isValid, nil
}

// func (s *NotificationService) GetNotificationDetail(ctx *gin.Context, id int64, notification *structs.NotificationDetail) (bool, bool, error) {
// 	userID := ctx.GetInt64("userID")

// 	if err := s.notificationRepository.GetNotificationByID2(ctx, id, notification); err != nil {
// 		return true, true, err
// 	}

// 	if notification.ID == 0 {
// 		return false, true, nil
// 	}

// 	if notification.SenderID != userID {
// 		return true, false, nil
// 	}

// 	return true, true, nil
// }
