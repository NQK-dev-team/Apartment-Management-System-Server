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
	userRepository          *repositories.UserRepository
	contractRepository      *repositories.ContractRepository
	supportTicketRepository *repositories.SupportTicketRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository:          repositories.NewUserRepository(),
		contractRepository:      repositories.NewContractRepository(),
		supportTicketRepository: repositories.NewSupportTicketRepository(),
	}
}

func (s *UserService) GetUsers(ctx *gin.Context, users *[]models.UserModel) error {
	if err := s.userRepository.Get(ctx, users); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(ctx *gin.Context, email string, user *models.UserModel) error {
	if err := s.userRepository.GetByEmail(ctx, user, email); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(ctx *gin.Context, id int64, user *models.UserModel) error {
	if err := s.userRepository.GetByID(ctx, user, id); err != nil {
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
		if err := s.userRepository.Create(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) UpdateUser(ctx *gin.Context, user *models.UserModel) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		oldData := models.UserModel{}
		if err := s.userRepository.GetByID(ctx, &oldData, user.ID); err != nil {
			return err
		}
		if oldData.Password != user.Password && !utils.CompareHashPassword(oldData.Password, user.Password) {
			hashedPassword, err := utils.HashPassword(user.Password)
			if err != nil {
				return err
			}
			user.Password = hashedPassword
		}
		if err := s.userRepository.Update(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) DeleteUsers(ctx *gin.Context, IDs []int64) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.DeleteByIDs(ctx, tx, IDs); err != nil {
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

		if err := s.userRepository.GetByIDs(ctx, users, []int64{claim.UserID}); err != nil {
			return err
		}
		return nil
	}

	if err := s.userRepository.GetStaffList(ctx, users); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetStaffDetail(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := s.userRepository.GetStaffDetail(ctx, user, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetStaffSchedule(ctx *gin.Context, schedules *[]models.ManagerScheduleModel, staffID int64) error {
	if err := s.userRepository.GetStaffSchedule(ctx, schedules, staffID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetStaffRelatedContract(ctx *gin.Context, contracts *[]models.ContractModel, staffID int64) error {
	if err := s.contractRepository.GetContractsByManagerID(ctx, contracts, staffID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetStaffRelatedTicket(ctx *gin.Context, tickets *[]structs.SupportTicket, staffID int64) error {
	if err := s.supportTicketRepository.GetTicketsByManagerID(ctx, tickets, staffID); err != nil {
		return err
	}
	return nil
}
