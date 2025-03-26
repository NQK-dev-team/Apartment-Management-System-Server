package repositories

import (
	"api/config"
	"api/models"
	"errors"

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

func (r *UserRepository) Create(ctx *gin.Context, user *models.UserModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx *gin.Context, user *models.UserModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx *gin.Context, users *models.UserModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := config.DB.Set("userID", userID).Model(&models.UserModel{}).Delete(users).Error; err != nil {
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
