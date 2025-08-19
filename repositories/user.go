package repositories

import (
	"api/config"
	"api/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetNewID(ctx *gin.Context) (int64, error) {
	lastestUser := models.UserModel{}
	if err := config.DB.Model(&models.UserModel{}).Order("id desc").Unscoped().Find(&lastestUser).Error; err != nil {
		return 0, err
	}
	return lastestUser.ID + 1, nil
}

func (r *UserRepository) GetByID(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := config.DB.Model(&models.UserModel{}).Where("id = ?", id).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByIDs(ctx *gin.Context, user *[]models.UserModel, id []int64) error {
	if err := config.DB.Model(&models.UserModel{}).Where("id IN ?", id).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetBySSN(ctx *gin.Context, user *models.UserModel, ssn string) error {
	if err := config.DB.Model(&models.UserModel{}).Where("ssn = ?", ssn).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByOldSSN(ctx *gin.Context, user *models.UserModel, ssn string) error {
	if err := config.DB.Model(&models.UserModel{}).Where("old_ssn = ?", ssn).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByPhone(ctx *gin.Context, user *models.UserModel, phone string) error {
	if err := config.DB.Model(&models.UserModel{}).Where("phone = ?", phone).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx *gin.Context, user *models.UserModel, email string) error {
	if err := config.DB.Model(&models.UserModel{}).Where("email = ?", email).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Get(ctx *gin.Context, user *[]models.UserModel) error {
	if err := config.DB.Model(&models.UserModel{}).Find(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Create(ctx *gin.Context, tx *gorm.DB, user *models.UserModel) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Omit("ID").Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx *gin.Context, tx *gorm.DB, user *models.UserModel, isQuiet bool) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}
	if err := tx.Set("userID", userID).Set("isQuiet", isQuiet).Save(user).Error; err != nil {
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
	if err := config.DB.Model(&models.UserModel{}).Where("id = ? AND is_owner = false AND is_manager = true AND is_customer = false", id).Find(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetStaffSchedule(ctx *gin.Context, schedules *[]models.ManagerScheduleModel, staffID int64) error {
	if err := config.DB.Model(&models.ManagerScheduleModel{}).Preload("Building").Preload("Manager").
		Joins("JOIN \"user\" ON \"user\".id = manager_id AND \"user\".deleted_at IS NULL").
		Where("manager_id = ?", staffID).Find(schedules).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetCustomerList(ctx *gin.Context, users *[]models.UserModel, limit int64, offset int64) error {
	if err := config.DB.Model(&models.UserModel{}).Where("is_owner = false AND is_manager = false AND is_customer = true").Limit(int(limit)).Offset(int(offset)).Find(users).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetCustomerDetail(ctx *gin.Context, user *models.UserModel, id int64) error {
	if err := config.DB.Model(&models.UserModel{}).Where("id = ? AND is_owner = false AND is_manager = false AND is_customer = true", id).Find(user).Error; err != nil {
		return err
	}
	return nil
}
