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
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

type BillService struct {
	billRepository     *repositories.BillRepository
	buildingService    *BuildingService
	contractService    *ContractService
	contractRepository *repositories.ContractRepository
}

func NewBillService() *BillService {
	return &BillService{
		billRepository:     repositories.NewBillRepository(),
		buildingService:    NewBuildingService(true),
		contractService:    NewContractService(),
		contractRepository: repositories.NewContractRepository(),
	}
}

func (s *BillService) GetBillList(ctx *gin.Context, bills *[]structs.Bill, limit, offset int64, startMonth, endMonth string) error {
	role := ctx.GetString("role")

	switch role {
	case constants.Roles.Manager, constants.Roles.Customer:
		if role == constants.Roles.Manager {
			if err := s.billRepository.GetBillListForManager(ctx, bills, startMonth, endMonth, limit, offset, ctx.GetInt64("userID")); err != nil {
				return err
			}
		} else {
			if err := s.billRepository.GetBillListForCustomer(ctx, bills, startMonth, endMonth, limit, offset, ctx.GetInt64("userID")); err != nil {
				return err
			}
		}

	case constants.Roles.Owner:
		if err := s.billRepository.GetBillList(ctx, bills, startMonth, endMonth, limit, offset); err != nil {
			return err
		}
	}

	return nil
}

func (s *BillService) DeleteBill(ctx *gin.Context, IDs []int64) (bool, error) {
	role := ctx.GetString("role")

	if constants.Roles.Owner == constants.Roles.Manager {
		userID := ctx.GetInt64("userID")

		bills := []models.BillModel{}
		if err := s.billRepository.GetDeletableBills(ctx, &bills, IDs, &userID); err != nil {
			return true, err
		}

		if len(bills) != len(IDs) {
			return false, nil
		}
	} else if role == constants.Roles.Owner {
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
	role := ctx.GetString("role")

	if role == constants.Roles.Manager || role == constants.Roles.Customer {
		if role == constants.Roles.Manager {
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
	bill := &structs.Bill{}
	if err := s.billRepository.GetBillByIDForCustomer(ctx, bill, ctx.GetInt64("userID"), billID); err != nil {
		return false
	}

	if bill.ID != billID {
		return false
	}

	return true
}

func (s *BillService) UpdateBill(ctx *gin.Context, bill *structs.UpdateBill, ID int64) (bool, bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager {
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
		return false, true, nil
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

	if bill.PayerID != 0 && bill.PaymentTime != "" {
		residentList := &[]models.RoomResidentModel{}

		if err := s.contractRepository.GetContractResidents(ctx, oldBill.ContractID, residentList); err != nil {
			return true, true, err
		}

		isPayerBelongToContract := false

		if oldBill.Contract.HouseholderID == bill.PayerID {
			isPayerBelongToContract = true
		}

		for _, resident := range *residentList {
			if resident.UserAccount.ID == bill.PayerID {
				isPayerBelongToContract = true
				break
			}
		}

		if !isPayerBelongToContract {
			return true, false, nil
		}
	} else if (bill.PayerID != 0 && bill.PaymentTime == "") || (bill.PayerID == 0 && bill.PaymentTime != "") {
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

	return true, true, err
}

func (s *BillService) AddBill(ctx *gin.Context, bill *structs.AddBill, newBillID *int64) (bool, bool, error) {
	role := ctx.GetString("role")

	if role == constants.Roles.Manager {
		isAllowed, err := s.contractService.CheckManagerContractPermission(ctx, ctx.GetInt64("userID"), bill.ContractID)
		if err != nil {
			return true, true, err
		}

		if !isAllowed {
			return false, true, nil
		}

	}

	contract := &structs.Contract{}
	if err := s.contractRepository.GetContractByID(ctx, contract, bill.ContractID); err != nil {
		return true, true, err
	}

	if contract.Status != constants.Common.ContractStatus.ACTIVE {
		return true, false, nil
	}

	if bill.Status == constants.Common.BillStatus.PAID {
		if bill.PayerID == 0 || bill.PaymentTime == "" {
			return true, false, nil
		}

		paymentTime, err := utils.ParseTimeWithZone(bill.PaymentTime + " 00:00:00")
		if err != nil {
			return true, true, err
		}

		if paymentTime.After(time.Now()) {
			return true, false, nil
		}

		periodTime, err := utils.ParseTimeWithZone(bill.Period + "-01 00:00:00")
		if err != nil {
			return true, true, err
		}

		if paymentTime.Before(periodTime) {
			return true, false, nil
		}
	}

	if bill.PayerID != 0 && bill.PaymentTime != "" {
		if bill.Status != constants.Common.BillStatus.PAID && bill.Status != constants.Common.BillStatus.CANCELLED {
			return true, false, nil
		}

		residentList := &[]models.RoomResidentModel{}

		if err := s.contractRepository.GetContractResidents(ctx, contract.ID, residentList); err != nil {
			return true, true, err
		}

		isPayerBelongToContract := false

		if contract.HouseholderID == bill.PayerID {
			isPayerBelongToContract = true
		}

		for _, resident := range *residentList {
			if resident.UserAccount.ID == bill.PayerID {
				isPayerBelongToContract = true
				break
			}
		}

		if !isPayerBelongToContract {
			return true, false, nil
		}
	} else if (bill.PayerID != 0 && bill.PaymentTime == "") || (bill.PayerID == 0 && bill.PaymentTime != "") {
		return true, false, nil
	}

	if len(bill.BillPayments) == 0 {
		return true, false, nil
	}

	billPeriod, err := utils.ParseTime(bill.Period + "-01")
	if err != nil {
		return true, true, err
	}

	newBill := &models.BillModel{
		Title:      bill.Title,
		Note:       sql.NullString{String: bill.Note, Valid: bill.Note != ""},
		Status:     bill.Status,
		ContractID: bill.ContractID,
		Period:     billPeriod,
		Amount:     0, // Will be calculated later
	}

	if bill.PayerID != 0 && bill.PaymentTime != "" {
		paymentTime, err := utils.ParseTimeWithZone(bill.PaymentTime + " 00:00:00")
		if err != nil {
			return true, true, err
		}

		newBill.PayerID = sql.NullInt64{
			Int64: bill.PayerID,
			Valid: bill.PayerID != 0,
		}
		newBill.PaymentTime = sql.NullTime{
			Time:  paymentTime,
			Valid: bill.PaymentTime != "",
		}
	} else {
		newBill.PayerID = sql.NullInt64{Valid: false}
		newBill.PaymentTime = sql.NullTime{Valid: false}
	}

	for _, payment := range bill.BillPayments {
		newBill.Amount += payment.Amount
		newBill.BillPayments = append(newBill.BillPayments, models.BillPaymentModel{
			Name:   payment.Name,
			Amount: payment.Amount,
			Note: sql.NullString{
				String: payment.Note,
				Valid:  payment.Note != "",
			},
		})
	}

	return true, true, config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.billRepository.CreateBill(ctx, tx, newBill); err != nil {
			return err
		}

		*newBillID = newBill.ID
		return nil
	})
}

func (s *BillService) UpdateBillStatus() error {
	return config.WorkerDB.Transaction(func(tx *gorm.DB) error {
		if err := s.billRepository.UpdateBillStatus(tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *BillService) InitBillPayment(ctx *gin.Context, billID int64, momoResponse *structs.MoMoCreatePaymentResponse, paymentResult *int) (bool, error) {
	if !s.CheckCustomerPermission(ctx, billID) {
		return false, nil
	}

	bill := &models.BillModel{}

	if err := s.billRepository.GetById2(ctx, bill, billID); err != nil {
		return true, err
	}

	if bill.Status != constants.Common.BillStatus.UN_PAID && bill.Status != constants.Common.BillStatus.OVERDUE {
		switch bill.Status {
		case constants.Common.BillStatus.PAID:
			*paymentResult = 1
		case constants.Common.BillStatus.PROCESSING:
			*paymentResult = 2
		case constants.Common.BillStatus.CANCELLED:
			return false, nil
		}
		return true, nil
	}

	flake, err := sonyflake.New(sonyflake.Settings{})
	if err != nil {
		fmt.Println(err)
		return true, errors.New("failed to create sonyflake")
	}

	requestID, _ := flake.NextID()
	orderID, _ := flake.NextID()

	return true, config.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.CreateMoMoPayment(bill, requestID, orderID, momoResponse); err != nil {
			return err
		}

		bill.RequestID = sql.NullString{
			String: fmt.Sprintf("%d", requestID),
			Valid:  true,
		}
		bill.OrderID = sql.NullString{
			String: fmt.Sprintf("%d", orderID),
			Valid:  true,
		}
		bill.Status = constants.Common.BillStatus.PROCESSING
		bill.PayerID = sql.NullInt64{
			Int64: ctx.GetInt64("userID"),
			Valid: true,
		}

		if err := s.billRepository.UpdateBill(ctx, tx, bill, bill.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *BillService) ProcessMoMoIPN(ctx *gin.Context, payload *structs.MoMoIPNPayload, billID int64) (bool, bool, error) {
	bill := &models.BillModel{}

	if err := s.billRepository.GetById2(ctx, bill, billID); err != nil {
		return true, true, err
	}

	if bill.ID != billID {
		return false, true, nil
	}

	if !utils.CheckIPNPayload(bill, payload) {
		return false, true, nil
	}

	var err error
	var isSuccess bool

	if payload.ResultCode == constants.Momo.ResultCode.Success || payload.ResultCode == constants.Momo.ResultCode.PaymentAuthorized {
		isSuccess = true

		err = config.DB.Transaction(func(tx *gorm.DB) error {
			bill.Status = constants.Common.BillStatus.PAID
			bill.PaymentTime = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}

			if err := s.billRepository.UpdateBill2(tx, bill, bill.ID); err != nil {
				return err
			}

			return nil
		})
	} else {
		isSuccess = false

		err = config.DB.Transaction(func(tx *gorm.DB) error {
			if err := s.billRepository.CancelBillPayment(tx, bill.ID); err != nil {
				return err
			}

			return nil
		})
	}

	return true, isSuccess, err
}

func (s *BillService) GetMomoResult() error {
	bills := []models.BillModel{}
	if err := s.billRepository.GetMomoBills(&bills); err != nil {
		return err
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for index, bill := range bills {
		tx.SavePoint(fmt.Sprintf("sp_%d", index))

		resultCode, err := utils.GetMoMoPaymentStatus(&bill)
		if err != nil {
			tx.RollbackTo(fmt.Sprintf("sp_%d", index))
			continue
		}

		if resultCode == constants.Momo.ResultCode.Success || resultCode == constants.Momo.ResultCode.PaymentAuthorized {
			bill.Status = constants.Common.BillStatus.PAID
			bill.PaymentTime = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}

			if err := s.billRepository.UpdateBill2(tx, &bill, bill.ID); err != nil {
				tx.RollbackTo(fmt.Sprintf("sp_%d", index))
				continue
			}
		} else if resultCode == constants.Momo.ResultCode.PaymentPending || resultCode == constants.Momo.ResultCode.PaymentProcessorPending || resultCode == constants.Momo.ResultCode.UserPending {
			continue
		} else {
			if err := s.billRepository.CancelBillPayment(tx, bill.ID); err != nil {
				tx.RollbackTo(fmt.Sprintf("sp_%d", index))
				continue
			}
		}
	}

	tx.Commit()

	return nil
}
