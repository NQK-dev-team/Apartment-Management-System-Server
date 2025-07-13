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

type BillService struct {
	billRepository *repositories.BillRepository
}

func NewBillService() *BillService {
	return &BillService{
		billRepository: repositories.NewBillRepository(),
	}
}

func (s *BillService) GetBillList(ctx *gin.Context, bills *[]structs.Bill, limit, offset int64, startMonth, endMonth string) error {
	role, exists := ctx.Get("role")
	if !exists {
		return errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager || role.(string) == constants.Roles.Customer {
		jwt, exists := ctx.Get("jwt")
		if !exists {
			return errors.New("jwt not found")
		}

		token, err := utils.ValidateJWTToken(jwt.(string))
		if err != nil {
			return err
		}

		claim := &structs.JTWClaim{}
		utils.ExtractJWTClaim(token, claim)

		if role.(string) == constants.Roles.Manager {
			if err := s.billRepository.GetBillListForManager(ctx, bills, startMonth, endMonth, limit, offset, claim.UserID); err != nil {
				return err
			}
		} else {
			if err := s.billRepository.GetBillListForCustomer(ctx, bills, startMonth, endMonth, limit, offset, claim.UserID); err != nil {
				return err
			}
		}

	} else if role.(string) == constants.Roles.Owner {
		if err := s.billRepository.GetBillList(ctx, bills, startMonth, endMonth, limit, offset); err != nil {
			return err
		}
	}

	return nil
}

func (s *BillService) DeleteBill(ctx *gin.Context, IDs []int64) (bool, error) {
	role, exists := ctx.Get("role")

	if !exists {
		return true, errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager {
		jwt, exists := ctx.Get("jwt")

		if !exists {
			return true, errors.New("jwt not found")
		}

		token, err := utils.ValidateJWTToken(jwt.(string))

		if err != nil {
			return true, err
		}

		claim := &structs.JTWClaim{}

		utils.ExtractJWTClaim(token, claim)

		bills := []models.BillModel{}
		if err := s.billRepository.GetDeletableBills(ctx, &bills, IDs, &claim.UserID); err != nil {
			return true, err
		}

		if len(bills) != len(IDs) {
			return false, nil
		}
	} else if role.(string) == constants.Roles.Owner {
		bills := []models.BillModel{}
		if err := s.billRepository.GetDeletableBills(ctx, &bills, IDs, nil); err != nil {
			return true, err
		}

		if len(bills) != len(IDs) {
			return false, nil
		}
	}
	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		return s.billRepository.Delete(ctx, tx, IDs)
	})
}
