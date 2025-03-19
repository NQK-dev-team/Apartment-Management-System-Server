package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"errors"

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
		if err := s.UserRepository.Create(ctx, tx, user); err != nil {
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
		if err := s.UserRepository.Update(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) DeleteUser(ctx *gin.Context, user *models.UserModel) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.UserRepository.Delete(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) GetStaffList(ctx *gin.Context, users *[]models.UserModel) error {
	role, exists := ctx.Get("role")

	if !exists {
		return errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager {
		jwt, exists := ctx.Get("jwt")

		if !exists {
			return errors.New("jwt not found")
		}

		token, err := utils.ValidateJWTToken(jwt.(string))

		if err != nil {
			return errors.New("jwt not valid")
		}

		claim := &structs.JTWClaim{}

		utils.ExtractJWTClaim(token, claim)

		if err := s.UserRepository.GetByIDs(ctx, users, []int64{claim.UserID}); err != nil {
			return err
		}
	}

	if err := s.UserRepository.GetStaffList(ctx, users); err != nil {
		return err
	}
	return nil
}
