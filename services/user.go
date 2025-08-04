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
	"time"

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

func (s *UserService) CreateStaff(ctx *gin.Context, newStaff *structs.NewStaff, newStaffID *int64) error {

	newPassword, err := utils.GeneratePassword(constants.Common.NewPasswordLength)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	profilePath := ""
	frontSSNPath := ""
	backSSNPath := ""

	err = config.DB.Transaction(func(tx *gorm.DB) error {
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
		}

		if err := s.userRepository.Create(ctx, tx, newUser); err != nil {
			return err
		}

		newID := newUser.ID

		newIDStr := strconv.FormatInt(newID, 10)

		profilePath, err = utils.StoreFile(newStaff.ProfileImage, constants.GetUserImageURL("images", newIDStr, ""))
		if err != nil {
			return err
		}

		frontSSNPath, err = utils.StoreFile(newStaff.FrontSSNImage, constants.GetUserImageURL("images", newIDStr, ""))
		if err != nil {
			return err
		}

		backSSNPath, err = utils.StoreFile(newStaff.BackSSNImage, constants.GetUserImageURL("images", newIDStr, ""))
		if err != nil {
			return err
		}

		newUser.ProfileFilePath = profilePath
		newUser.SSNFrontFilePath = frontSSNPath
		newUser.SSNBackFilePath = backSSNPath

		if err := s.userRepository.Update(ctx, tx, newUser, true); err != nil {
			return err
		}

		schedules := []models.ManagerScheduleModel{}
		for _, val := range newStaff.Schedules {
			startDate, err := utils.ParseTime(val.StartDate)
			if err != nil {
				return err
			}

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

		*newStaffID = newUser.ID

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
				startDate, err := utils.ParseTime(schedule.StartDate)
				if err != nil {
					return err
				}
				schedules = append(schedules, models.ManagerScheduleModel{
					BuildingID: schedule.BuildingID,
					ManagerID:  editStaff.ID,
					StartDate:  startDate,
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
						startDate, err := utils.ParseTime(val.StartDate)
						if err != nil {
							return err
						}
						schedules[index].BuildingID = val.BuildingID
						schedules[index].StartDate = startDate
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
		if err := s.userRepository.Update(ctx, tx, user, false); err != nil {
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

func (s *UserService) GetStaffRelatedTicket(ctx *gin.Context, tickets *[]structs.SupportTicket, staffID int64, limit int64, offset int64, startDate string, endDate string) error {
	if err := s.supportTicketRepository.GetTicketsByManagerID(ctx, tickets, staffID, limit, offset, startDate, endDate); err != nil {
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

func (s *UserService) CheckDuplicateData2(ctx *gin.Context, ssn string, phone string, oldSSN string) (string, error) {
	jwt, exists := ctx.Get("jwt")

	if !exists {
		return "", errors.New("jwt not found")
	}

	token, err := utils.ValidateJWTToken(jwt.(string))

	if err != nil {
		return "", errors.New("jwt not valid")
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	user := &models.UserModel{}

	if err := s.userRepository.GetBySSN(ctx, user, ssn); err != nil {
		return "", err
	}

	if user.ID != 0 && user.ID != claim.UserID {
		return config.GetMessageCode("SSN_ALREADY_EXISTS"), nil
	}

	user = &models.UserModel{}

	if err := s.userRepository.GetByPhone(ctx, user, phone); err != nil {
		return "", err
	}

	if user.ID != 0 && user.ID != claim.UserID {
		return config.GetMessageCode("PHONE_ALREADY_EXISTS"), nil
	}

	user = &models.UserModel{}

	if oldSSN != "" {
		if err := s.userRepository.GetByOldSSN(ctx, user, oldSSN); err != nil {
			return "", err
		}

		if user.ID != 0 && user.ID != claim.UserID {
			return config.GetMessageCode("OLD_SSN_ALREADY_EXISTS"), nil
		}
	}

	return "", nil
}

func (s *UserService) GetCustomerList(ctx *gin.Context, users *[]models.UserModel, limit int64, offset int64) error {
	if err := s.userRepository.GetCustomerList(ctx, users, limit, offset); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetCustomerDetail(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := s.userRepository.GetCustomerDetail(ctx, user, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetCustomerContract(ctx *gin.Context, contracts *[]structs.Contract, customerID int64) error {
	if err := s.contractRepository.GetContractsByCustomerID(ctx, contracts, customerID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetCustomerTicket(ctx *gin.Context, tickets *[]structs.SupportTicket, customerID int64) error {
	if err := s.supportTicketRepository.GetTicketsByCustomerID(ctx, tickets, customerID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) CreateCustomer(ctx *gin.Context, newCustomer *structs.NewCustomer, newCustomerID *int64) error {
	newPassword, err := utils.GeneratePassword(constants.Common.NewPasswordLength)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	profilePath := ""
	frontSSNPath := ""
	backSSNPath := ""

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		newUser := &models.UserModel{
			LastName:  newCustomer.LastName,
			FirstName: newCustomer.FirstName,
			MiddleName: sql.NullString{
				String: newCustomer.MiddleName,
				Valid:  newCustomer.MiddleName != "",
			},
			DOB:    newCustomer.Dob,
			POB:    newCustomer.Pob,
			Gender: newCustomer.Gender,
			SSN:    newCustomer.SSN,
			OldSSN: sql.NullString{
				String: newCustomer.OldSSN,
				Valid:  newCustomer.OldSSN != "",
			},
			Email:            newCustomer.Email,
			Phone:            newCustomer.Phone,
			PermanentAddress: newCustomer.PermanentAddress,
			TemporaryAddress: newCustomer.TemporaryAddress,
			Password:         hashedPassword,
			IsOwner:          false,
			IsManager:        false,
			IsCustomer:       true,
		}

		if err := s.userRepository.Create(ctx, tx, newUser); err != nil {
			return err
		}

		newID := newUser.ID
		newIDStr := strconv.FormatInt(newID, 10)

		profilePath, err = utils.StoreFile(newCustomer.ProfileImage, constants.GetUserImageURL("images", newIDStr, ""))
		if err != nil {
			return err
		}

		frontSSNPath, err = utils.StoreFile(newCustomer.FrontSSNImage, constants.GetUserImageURL("images", newIDStr, ""))
		if err != nil {
			return err
		}

		backSSNPath, err = utils.StoreFile(newCustomer.BackSSNImage, constants.GetUserImageURL("images", newIDStr, ""))
		if err != nil {
			return err
		}
		newUser.ProfileFilePath = profilePath
		newUser.SSNFrontFilePath = frontSSNPath
		newUser.SSNBackFilePath = backSSNPath

		if err := s.userRepository.Update(ctx, tx, newUser, true); err != nil {
			return err
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

		*newCustomerID = newUser.ID

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

func (s *UserService) GetUserInfo(ctx *gin.Context, user *models.UserModel) error {
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

	if err := s.userRepository.GetByID(ctx, user, claim.UserID); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateProfile(ctx *gin.Context, profile *structs.UpdateProfile) error {
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

	deleteFileList := []string{}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		user := &models.UserModel{}

		if err := s.userRepository.GetByID(ctx, user, claim.UserID); err != nil {
			return err
		}

		user.FirstName = profile.FirstName
		user.LastName = profile.LastName
		user.Phone = profile.Phone
		user.PermanentAddress = profile.PermanentAddress
		user.TemporaryAddress = profile.TemporaryAddress
		user.MiddleName = sql.NullString{
			String: profile.MiddleName,
			Valid:  profile.MiddleName != "",
		}
		user.DOB = profile.Dob
		user.POB = profile.Pob
		user.Gender = profile.Gender
		user.SSN = profile.SSN
		user.OldSSN = sql.NullString{
			String: profile.OldSSN,
			Valid:  profile.OldSSN != "",
		}

		IDStr := strconv.FormatInt(user.ID, 10)

		if profile.NewProfileImage != nil {
			profilePath, err := utils.StoreFile(profile.NewProfileImage, constants.GetUserImageURL("images", IDStr, ""))
			if err != nil {
				return err
			}
			deleteFileList = append(deleteFileList, user.ProfileFilePath)
			user.ProfileFilePath = profilePath
		}

		if profile.NewFrontSSNImage != nil {
			frontSSNPath, err := utils.StoreFile(profile.NewFrontSSNImage, constants.GetUserImageURL("images", IDStr, ""))
			if err != nil {
				return err
			}
			deleteFileList = append(deleteFileList, user.SSNFrontFilePath)
			user.SSNFrontFilePath = frontSSNPath
		}

		if profile.NewBackSSNImage != nil {
			backSSNPath, err := utils.StoreFile(profile.NewBackSSNImage, constants.GetUserImageURL("images", IDStr, ""))
			if err != nil {
				return err
			}
			deleteFileList = append(deleteFileList, user.SSNBackFilePath)
			user.SSNBackFilePath = backSSNPath
		}

		if err := s.userRepository.Update(ctx, tx, user, false); err != nil {
			return err
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

func (s *UserService) ChangePassword(ctx *gin.Context, changePassword *structs.ChangePassword) (bool, error) {
	jwt, exists := ctx.Get("jwt")

	if !exists {
		return true, errors.New("jwt not found")
	}

	token, err := utils.ValidateJWTToken(jwt.(string))

	if err != nil {
		return true, errors.New("jwt not valid")
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	user := &models.UserModel{}

	if err := s.userRepository.GetByID(ctx, user, claim.UserID); err != nil {
		return true, err
	}

	if !utils.CompareHashPassword(user.Password, changePassword.OldPassword) {
		return false, nil
	}

	if changePassword.OldPassword == changePassword.NewPassword {
		return true, nil
	}

	hashedPassword, err := utils.HashPassword(changePassword.NewPassword)
	if err != nil {
		return true, err
	}

	user.Password = hashedPassword

	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.Update(ctx, config.DB, user, false); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserService) ChangeEmail(ctx *gin.Context, changeEmail *structs.ChangeEmail) (bool, bool, error) {
	jwt, exists := ctx.Get("jwt")
	if !exists {
		return true, true, errors.New("jwt not found")
	}

	token, err := utils.ValidateJWTToken(jwt.(string))
	if err != nil {
		return true, true, errors.New("jwt not valid")
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	user := &models.UserModel{}
	if err := s.userRepository.GetByEmail(ctx, user, changeEmail.NewEmail); err != nil {
		return true, true, err
	}

	if user.ID != 0 && user.ID != claim.UserID {
		return true, false, nil
	}

	user = &models.UserModel{}

	if err := s.userRepository.GetByID(ctx, user, claim.UserID); err != nil {
		return true, true, err
	}

	if !utils.CompareHashPassword(user.Password, changeEmail.Password) {
		return false, true, nil
	}

	return true, true, config.DB.Transaction(func(tx *gorm.DB) error {
		user.Email = changeEmail.NewEmail
		user.EmailVerifiedAt = sql.NullTime{
			Valid: false,
			Time:  time.Time{},
		}

		if err := s.userRepository.Update(ctx, tx, user, false); err != nil {
			return err
		}

		if _, err := s.emailService.SendAccountChangeEmailVerificationEmail(ctx, changeEmail.NewEmail); err != nil {
			return err
		}

		return nil
	})
}
