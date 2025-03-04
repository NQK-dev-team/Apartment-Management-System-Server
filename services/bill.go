package services

import (
	"api/models"
	"api/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type BillService struct {
	billRepository *repositories.BillRepository
}

func NewBillService() *BillService {
	return &BillService{
		billRepository: repositories.NewBillRepository(),
	}
}

func (s *BillService) DeleteWithoutTransaction(ctx *gin.Context, id int64) error {
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = 0
	}

	now := time.Now()

	bill := &models.BillModel{}
	if err := s.billRepository.GetById(ctx, bill, id); err != nil {
		return err
	}

	bill.DefaultModel.DeletedBy = userID.(int64)
	bill.DefaultModel.DeletedAt.Valid = true
	bill.DefaultModel.DeletedAt.Time = now

	for index := range bill.ExtraPayments {
		bill.ExtraPayments[0].DefaultModel.DeletedBy = userID.(int64)
		bill.ExtraPayments[index].DefaultModel.DeletedAt.Valid = true
		bill.ExtraPayments[index].DefaultModel.DeletedAt.Time = now
	}

	if err := s.billRepository.QuietUpdate(ctx, bill); err != nil {
		return err
	}

	return nil
}
