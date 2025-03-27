package repositories

import (
	"api/config"
	"api/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByID(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := config.DB.Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetByIDs(ctx *gin.Context, user *[]models.UserModel, id []int64) error {
	if err := config.DB.Where("id = ?", id).Find(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetBySSN(ctx *gin.Context, user *models.UserModel, ssn string) error {
	if err := config.DB.Where("ssn = ?", ssn).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx *gin.Context, user *models.UserModel, email string) error {
	if err := config.DB.Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *UserRepository) Get(ctx *gin.Context, user *[]models.UserModel) error {
	if err := config.DB.Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Create(ctx *gin.Context, tx *gorm.DB, user *models.UserModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx *gin.Context, tx *gorm.DB, user *models.UserModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteByIDs(ctx *gin.Context, tx *gorm.DB, ids []int64) error {
	now := time.Now()
	userID := ctx.GetInt64("userID")

	if err := tx.Set("isQuiet", true).Model(&models.UserModel{}).Where("id in ?", ids).UpdateColumns(models.UserModel{
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

func (r *UserRepository) GetStaffList(ctx *gin.Context, users *[]models.UserModel) error {
	if err := config.DB.Model(&models.UserModel{}).Where("is_owner = false AND is_manager = true AND is_customer = false").Find(users).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetStaffDetail(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := config.DB.Where("id = ? AND is_owner = false AND is_manager = true AND is_customer = false", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetStaffSchedule(ctx *gin.Context, schedules *[]models.ManagerScheduleModel, staffID int64) error {
	if err := config.DB.Preload("Building").Preload("Manager").Where("manager_id = ?", staffID).Find(schedules).Error; err != nil {
		return err
	}
	return nil
}
