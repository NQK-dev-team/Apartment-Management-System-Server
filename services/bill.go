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
	billRepository  *repositories.BillRepository
	buildingService *BuildingService
}

func NewBillService() *BillService {
	return &BillService{
		billRepository:  repositories.NewBillRepository(),
		buildingService: NewBuildingService(true),
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

func (s *BillService) GetBillDetail(ctx *gin.Context, bill *structs.Bill, billID int64) (bool, error) {
	role, exists := ctx.Get("role")
	if !exists {
		return true, errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager || role.(string) == constants.Roles.Customer {
		if role.(string) == constants.Roles.Manager {
			if isAllowed := s.CheckManagerPermission(ctx, billID); !isAllowed {
				return false, nil
			}
		} else {
			if isAllowed := s.CheckCustomerPermission(ctx, billID); !isAllowed {
				return false, nil
			}
		}
	}

	if err := s.billRepository.GetById(ctx, bill, billID); err != nil {
		return true, err
	}

	return true, nil
}

func (s *BillService) CheckManagerPermission(ctx *gin.Context, billID int64) bool {
	var buildingID int64
	if err := s.billRepository.GetBillBuildingID(ctx, billID, &buildingID); err != nil {
		return false
	}

	return s.buildingService.CheckManagerPermission(ctx, buildingID)
}

func (s *BillService) CheckCustomerPermission(ctx *gin.Context, billID int64) bool {
	jwt, exists := ctx.Get("jwt")
	if !exists {
		return false
	}

	token, err := utils.ValidateJWTToken(jwt.(string))
	if err != nil {
		return false
	}

	claim := &structs.JTWClaim{}
	utils.ExtractJWTClaim(token, claim)

	bill := &structs.Bill{}
	if err := s.billRepository.GetBillByIDForCustomer(ctx, bill, claim.UserID, billID); err != nil {
		return false
	}

	if bill.ID != billID {
		return false
	}

	return true
}
