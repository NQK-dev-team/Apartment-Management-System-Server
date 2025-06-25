package repositories

import (
	"api/config"
	"api/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomRepository struct{}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (r *RoomRepository) GetById(ctx *gin.Context, room *models.RoomModel, id int64) error {
	if err := config.DB.Model(&models.RoomModel{}).
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("room.id = ? AND room.deleted_at IS NULL", id).Preload("Images").Preload("Contracts").Find(room).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil
		// }
		return err
	}
	return nil
}

func (r *RoomRepository) GetByIDs(ctx *gin.Context, rooms *[]models.RoomModel, IDs []int64) error {
	if err := config.DB.Model(&models.RoomModel{}).
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("room.id in ? AND room.deleted_at IS NULL", IDs). // Preload("Images").Preload("Contracts").
		Find(rooms).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil
		// }
		return err
	}
	return nil
}

func (r *RoomRepository) GetNewID(ctx *gin.Context) (int64, error) {
	lastestRoom := models.RoomModel{}
	if err := config.DB.Model(&models.RoomModel{}).Order("id desc").Unscoped().Find(&lastestRoom).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return 0, nil
		// }
		return 0, err
	}
	return lastestRoom.ID + 1, nil
}

func (r *RoomRepository) GetNewImageID(ctx *gin.Context) (int64, error) {
	lastestImage := models.RoomImageModel{}
	if err := config.DB.Model(&models.RoomImageModel{}).Order("id desc").Unscoped().Find(&lastestImage).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return 0, nil
		// }
		return 0, err
	}
	return lastestImage.ID + 1, nil
}

func (r *RoomRepository) Delete(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.RoomModel{}).Where("id in ?", id).UpdateColumns(models.RoomModel{
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

	if err := tx.Set("isQuiet", true).Model(&models.RoomImageModel{}).Where("room_id in ?", id).UpdateColumns(models.RoomImageModel{
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

func (r *RoomRepository) Create(ctx *gin.Context, tx *gorm.DB, rooms *[]models.RoomModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Model(&models.RoomModel{}).Omit("ID").Create(rooms).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) CreateImage(ctx *gin.Context, tx *gorm.DB, images *[]models.RoomImageModel) error {
	userID := ctx.GetInt64("userID")
	if err := tx.Set("userID", userID).Model(&models.RoomImageModel{}).Omit("ID").Create(images).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) DeleteImages(ctx *gin.Context, tx *gorm.DB, id []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.RoomImageModel{}).Where("id in ?", id).UpdateColumns(models.RoomImageModel{
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

func (r *RoomRepository) Update(ctx *gin.Context, tx *gorm.DB, rooms *[]models.RoomModel) error {
	userID := ctx.GetInt64("userID")
	// if err := tx.Set("userID", userID).Save(rooms).Error; err != nil {
	// 	return err
	// }

	for _, room := range *rooms {
		if err := tx.Set("userID", userID).Model(&models.RoomModel{}).Where("id = ?", room.ID).Updates(room).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *RoomRepository) GetNewImageNo(ctx *gin.Context, roomID int64) (int, error) {
	lastestImage := models.RoomImageModel{}
	if err := config.DB.Model(&models.RoomImageModel{}).Where("room_id = ?", roomID).Order("no desc").Unscoped().Find(&lastestImage).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return 0, nil
		// }
		return 0, err
	}
	return lastestImage.No + 1, nil
}

func (r *RoomRepository) GetRoomByRoomIDAndBuildingID(ctx *gin.Context, room *models.RoomModel, roomID int64, buildingID int64) error {
	if err := config.DB.Model(&models.RoomModel{}).
		Joins("JOIN building ON building.id = room.building_id AND building.deleted_at IS NULL").
		Where("room.id = ? AND room.building_id = ? AND room.deleted_at IS NULL", roomID, buildingID).Preload("Images").Preload("Contracts").Find(room).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil
		// }
		return err
	}
	return nil
}
