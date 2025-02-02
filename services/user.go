package services

import (
	"api/config"
	"api/models"
	"api/repositories"
	"api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService() *UserService {
	userRepository := repositories.NewUserRepository()
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) GetUsers(ctx *gin.Context, users *[]models.UserModel) error {
	if err := s.UserRepository.Get(ctx, users); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(ctx *gin.Context, email string, user *models.UserModel) error {
	if err := s.UserRepository.GetByEmail(ctx, user, email); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(ctx *gin.Context, id int64, user *models.UserModel) error {
	if err := s.UserRepository.GetByID(ctx, user, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) CreateUser(ctx *gin.Context, user *models.UserModel) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
		if err := s.UserRepository.Create(ctx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) UpdateUser(ctx *gin.Context, user *models.UserModel) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		oldData := models.UserModel{}
		if err := s.UserRepository.GetByID(ctx, &oldData, user.ID); err != nil {
			return err
		}
		if oldData.Password != user.Password && !utils.CompareHashPassword(oldData.Password, user.Password) {
			hashedPassword, err := utils.HashPassword(user.Password)
			if err != nil {
				return err
			}
			user.Password = hashedPassword
		}
		if err := s.UserRepository.Update(ctx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) DeleteUser(ctx *gin.Context, user *models.UserModel) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.UserRepository.Delete(ctx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}
