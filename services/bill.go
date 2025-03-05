package services

import (
	"api/models"
	"api/repositories"

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

func (s *BillService) GetBill(ctx *gin.Context, bill *[]models.BillModel) (bool, error) {
	// role, exists := ctx.Get("role")

	// if !exists {
	// 	return false, nil
	// }

	// if role.(string) == constants.Roles.Manager {
	// 	jwt, exists := ctx.Get("jwt")

	// 	if !exists {
	// 		return false, nil
	// 	}

	// 	token, err := utils.ValidateJWTToken(jwt.(string))

	// 	if err != nil {
	// 		return true, err
	// 	}

	// 	claim := &structs.JTWClaim{}

	// 	utils.ExtractJWTClaim(token, claim)

	// 	return true, s.billRepository.GetBillBaseOnSchedule(ctx, bill, claim.UserID)
	// }

	return true, s.billRepository.Get(ctx, bill)
}
