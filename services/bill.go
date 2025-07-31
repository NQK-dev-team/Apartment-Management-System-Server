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
	"time"

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

func (s *BillService) UpdateBill(ctx *gin.Context, bill *structs.UpdateBill, ID int64) (bool, bool, error) {
	role, exists := ctx.Get("role")
	if !exists {
		return true, true, errors.New("role not found")
	}

	if role.(string) == constants.Roles.Manager {
		if isAllowed := s.CheckManagerPermission(ctx, ID); !isAllowed {
			return false, true, nil
		}
	}

	oldBill := &models.BillModel{}
	if err := s.billRepository.GetById2(ctx, oldBill, ID); err != nil {
		return true, true, err
	}

	if oldBill.ID == 0 {
		return true, true, errors.New("bill not found")
	}

	if oldBill.Status == constants.Common.BillStatus.PAID || oldBill.Status == constants.Common.BillStatus.PROCESSING {
		return true, false, nil
	}

	if oldBill.Status != bill.Status && bill.Status != constants.Common.BillStatus.CANCELLED && bill.Status != constants.Common.BillStatus.PAID {
		return true, false, nil
	}

	if bill.PaymentTime != "" {
		payTime, err := utils.ParseTimeWithZone(bill.PaymentTime + " 00:00:00")
		if err != nil {
			return true, true, err
		}

		if payTime.After(time.Now()) {
			return true, false, nil
		}

		periodTime, err := utils.ParseTimeWithZone(oldBill.Period.Format("2006-01-02 15:04:05"))
		if err != nil {
			return true, true, err
		}

		if payTime.Before(periodTime) {
			return true, false, nil
		}
	}

	if len(bill.NewPayments) == 0 && len(bill.Payments) == 0 {
		return true, false, nil
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if len(bill.DeletedPayments) > 0 {
			if err := s.billRepository.DeletePayment(ctx, tx, bill.DeletedPayments); err != nil {
				return err
			}
		}

		if len(bill.NewPayments) > 0 {
			newPaymentModels := []models.BillPaymentModel{}
			for _, payment := range bill.NewPayments {
				newPaymentModels = append(newPaymentModels, models.BillPaymentModel{
					BillID: ID,
					Name:   payment.Name,
					Amount: payment.Amount,
					Note: sql.NullString{
						String: payment.Note,
						Valid:  payment.Note != "",
					},
				})
			}

			if err := s.billRepository.AddNewPayment(ctx, tx, &newPaymentModels); err != nil {
				return err
			}
		}

		if len(bill.Payments) > 0 {
			paymentModels := []models.BillPaymentModel{}
			for _, payment := range bill.Payments {
				paymentModels = append(paymentModels, models.BillPaymentModel{
					BillID: ID,
					Name:   payment.Name,
					Amount: payment.Amount,
					Note: sql.NullString{
						String: payment.Note,
						Valid:  payment.Note != "",
					},
					DefaultModel: models.DefaultModel{
						ID: payment.ID,
					},
				})
			}

			if err := s.billRepository.UpdatePayment(ctx, tx, &paymentModels); err != nil {
				return err
			}
		}

		var totalAmount float64 = 0
		for _, payment := range bill.NewPayments {
			totalAmount += payment.Amount
		}
		for _, payment := range bill.Payments {
			totalAmount += payment.Amount
		}

		oldBill.Title = bill.Title
		oldBill.Note = sql.NullString{
			String: bill.Note,
			Valid:  bill.Note != "",
		}
		oldBill.Amount = totalAmount
		if oldBill.Status != constants.Common.BillStatus.PAID && bill.Status == constants.Common.BillStatus.PAID && bill.PayerID != 0 && bill.PaymentTime != "" {
			oldBill.PayerID = sql.NullInt64{
				Int64: bill.PayerID,
				Valid: bill.PayerID != 0,
			}

			paymentTime, err := utils.ParseTimeWithZone(bill.PaymentTime + " 00:00:00")
			if err != nil {
				return err
			}

			oldBill.PaymentTime = sql.NullTime{
				Time:  paymentTime,
				Valid: bill.PaymentTime != "",
			}
		}
		oldBill.Status = bill.Status

		if err := s.billRepository.UpdateBill(ctx, tx, oldBill, ID); err != nil {
			return err
		}

		return nil
	})

	return true, err
	return true, true, err
}
