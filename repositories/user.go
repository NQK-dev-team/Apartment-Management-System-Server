package repositories

import (
	"api/config"
	"api/models"

	"github.com/gin-gonic/gin"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByID(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := config.DB.Where("id = ?", id).First(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetBySSN(ctx *gin.Context, user *models.UserModel, ssn string) error {
	if err := config.DB.Where("ssn = ?", ssn).First(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx *gin.Context, user *models.UserModel, email string) error {
	if err := config.DB.Where("email = ?", email).First(user).Error; err != nil {
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
	if err := config.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx *gin.Context, user *models.UserModel) error {
	if err := config.DB.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx *gin.Context, users *models.UserModel) error {
	if err := config.DB.Model(&models.UserModel{}).Delete(users).Error; err != nil {
		return err
	}
	return nil
}
