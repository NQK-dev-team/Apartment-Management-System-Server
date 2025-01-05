package services

import (
	"api/models"
	"api/repositories"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService() *UserService {
	userRepository := repositories.NewUserRepository()
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) GetUser(ctx *gin.Context) (*[]models.UserModel, error) {
	users := &[]models.UserModel{}
	return users, nil
}

func (s *UserService) CreateUser(ctx *gin.Context, users *[]models.UserModel) error {
	return nil
}

func (s *UserService) UpdateUser(ctx *gin.Context, users *[]models.UserModel) error {
	return nil
}

func (s *UserService) DeleteUser(ctx *gin.Context, users *[]models.UserModel) error {
	return nil
}
