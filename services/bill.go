package services

import (
	"api/repositories"
)

type BillService struct {
	billRepository *repositories.BillRepository
}

func NewBillService() *BillService {
	return &BillService{
		billRepository: repositories.NewBillRepository(),
	}
}
