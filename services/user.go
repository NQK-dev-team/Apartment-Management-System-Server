package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"database/sql"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	userRepository            *repositories.UserRepository
	contractRepository        *repositories.ContractRepository
	supportTicketRepository   *repositories.SupportTicketRepository
	managerScheduleRepository *repositories.ManagerScheduleRepository
	emailService              *EmailService
}

func NewUserService() *UserService {
	return &UserService{
		userRepository:            repositories.NewUserRepository(),
		contractRepository:        repositories.NewContractRepository(),
		supportTicketRepository:   repositories.NewSupportTicketRepository(),
		managerScheduleRepository: repositories.NewManagerScheduleRepository(),
		emailService:              NewEmailService(),
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

// func (s *UserService) CreateUser(ctx *gin.Context, user *models.UserModel) error {
// 	err := config.DB.Transaction(func(tx *gorm.DB) error {
// 		hashedPassword, err := utils.HashPassword(user.Password)
// 		if err != nil {
// 			return err
// 		}
// 		user.Password = hashedPassword
// 		if err := s.userRepository.Create(ctx, tx, user); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return err
// }

func (s *UserService) CreateStaff(ctx *gin.Context, newStaff *structs.NewStaff) error {

	newPassword, err := utils.GeneratePassword(8)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	newID, err := s.userRepository.GetNewID(ctx)

	if err != nil {
		return err
	}

	newIDStr := strconv.FormatInt(newID, 10)

	profilePath, err := utils.StoreFile(newStaff.ProfileImage, "images/users/"+newIDStr+"/")
	if err != nil {
		return err
	}

	frontSSNPath, err := utils.StoreFile(newStaff.FrontSSNImage, "images/users/"+newIDStr+"/")
	if err != nil {
		return err
	}

	backSSNPath, err := utils.StoreFile(newStaff.BackSSNImage, "images/users/"+newIDStr+"/")
	if err != nil {
		return err
	}

	newUser := &models.UserModel{
		LastName:  newStaff.LastName,
		FirstName: newStaff.FirstName,
		MiddleName: sql.NullString{
			String: newStaff.MiddleName,
			Valid:  newStaff.MiddleName != "",
		},
		DOB:    newStaff.Dob,
		POB:    newStaff.Pob,
		Gender: newStaff.Gender,
		SSN:    newStaff.SSN,
		OldSSN: sql.NullString{
			String: newStaff.OldSSN,
			Valid:  newStaff.OldSSN != "",
		},
		Email:            newStaff.Email,
		Phone:            newStaff.Phone,
		PermanentAddress: newStaff.PermanentAddress,
		TemporaryAddress: newStaff.TemporaryAddress,
		Password:         hashedPassword,
		IsOwner:          false,
		IsManager:        true,
		IsCustomer:       false,
		ProfileFilePath:  profilePath,
		SSNFrontFilePath: frontSSNPath,
		SSNBackFilePath:  backSSNPath,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.Create(ctx, tx, newUser); err != nil {
			return err
		}

		schedules := []models.ManagerScheduleModel{}
		for _, val := range newStaff.Schedules {
			startDate := utils.ParseTime(val.StartDate)
			endDate := utils.StringToNullTime(val.EndDate)
			schedules = append(schedules, models.ManagerScheduleModel{
				BuildingID: val.BuildingID,
				ManagerID:  newID,
				StartDate:  startDate,
				EndDate:    endDate,
			})
		}
		if len(schedules) > 0 {
			if err := s.managerScheduleRepository.Create(ctx, tx, &schedules); err != nil {
				return err
			}
		}

		fullName := ""

		if newUser.MiddleName.Valid {
			fullName = newUser.LastName + " " + newUser.MiddleName.String + " " + newUser.FirstName
		} else {
			fullName = newUser.LastName + " " + newUser.FirstName
		}

		if err := s.emailService.SendAccountCreationEmail(ctx, newUser.Email, fullName, newPassword); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.RemoveFile(profilePath)
		utils.RemoveFile(frontSSNPath)
		utils.RemoveFile(backSSNPath)
		return err
	}

	return nil
}

func (s *UserService) UpdateStaff(ctx *gin.Context, editStaff *structs.EditStaff) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if len(editStaff.DeletedSchedules) > 0 {
			if err := s.managerScheduleRepository.Delete(ctx, tx, editStaff.DeletedSchedules); err != nil {
				return err
			}
		}

		if len(editStaff.NewSchedules) > 0 {
			schedules := []models.ManagerScheduleModel{}
			for _, schedule := range editStaff.NewSchedules {
				schedules = append(schedules, models.ManagerScheduleModel{
					BuildingID: schedule.BuildingID,
					ManagerID:  editStaff.ID,
					StartDate:  utils.ParseTime(schedule.StartDate),
					EndDate:    utils.StringToNullTime(schedule.EndDate),
				})
			}
			if err := s.managerScheduleRepository.Create(ctx, tx, &schedules); err != nil {
				return err
			}
		}

		{
			scheduleIDs := []int64{}
			for _, schedule := range editStaff.Schedules {
				scheduleIDs = append(scheduleIDs, schedule.ID)
			}
			schedules := []models.ManagerScheduleModel{}
			if err := s.managerScheduleRepository.GetByIDs(ctx, &schedules, scheduleIDs); err != nil {
				return err
			}
			for index, schedule := range schedules {
				for _, val := range editStaff.Schedules {
					if val.ID == schedule.ID {
						schedules[index].BuildingID = val.BuildingID
						schedules[index].StartDate = utils.ParseTime(val.StartDate)
						schedules[index].EndDate = utils.StringToNullTime(val.EndDate)
						break
					}
				}
			}
			if err := s.managerScheduleRepository.Update(ctx, tx, &schedules); err != nil {
				return err
			}
		}
		return nil
	})
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

func (s *UserService) GetStaffRelatedTicket(ctx *gin.Context, tickets *[]structs.SupportTicket, staffID int64, limit int64, offset int64, quarters []struct {
	Year     int
	Quarters []int
}) error {
	if err := s.supportTicketRepository.GetTicketsByManagerID(ctx, tickets, staffID, limit, offset, quarters); err != nil {
		return err
	}
	return nil
}

func (s *UserService) CheckDuplicateData(ctx *gin.Context, email string, ssn string, phone string, oldSSN string) (string, error) {
	user := &models.UserModel{}

	if err := s.userRepository.GetByEmail(ctx, user, email); err != nil {
		return "", err
	}

	if user.ID != 0 {
		return config.GetMessageCode("EMAIL_ALREADY_EXISTS"), nil
	}

	if err := s.userRepository.GetBySSN(ctx, user, ssn); err != nil {
		return "", err
	}

	if user.ID != 0 {
		return config.GetMessageCode("SSN_ALREADY_EXISTS"), nil
	}

	if err := s.userRepository.GetByPhone(ctx, user, phone); err != nil {
		return "", err
	}

	if user.ID != 0 {
		return config.GetMessageCode("PHONE_ALREADY_EXISTS"), nil
	}

	if oldSSN != "" {
		if err := s.userRepository.GetByOldSSN(ctx, user, oldSSN); err != nil {
			return "", err
		}

		if user.ID != 0 {
			return config.GetMessageCode("OLD_SSN_ALREADY_EXISTS"), nil
		}
	}

	return "", nil
}
