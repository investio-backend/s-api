package service

import (
	"errors"

	"gitlab.com/investio/backend/sim-api/db"
	"gitlab.com/investio/backend/sim-api/v1/model"
	"gorm.io/gorm"
)

type TransactionService interface {
	Get(transList *[]model.Transaction, userID uint) (err error)
	Write(tran *model.Transaction) (err error)
}

type transactionService struct {
}

func NewTransctionService() TransactionService {
	return &transactionService{}
}

func (s *transactionService) Get(transList *[]model.Transaction, userID uint) (err error) {
	if err = db.SimDB.Limit(50).Where("user_id = ?", userID).Order("data_date desc").Find(transList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not found")
		}
	}
	return
}

func (s *transactionService) Write(tran *model.Transaction) (err error) {
	err = db.SimDB.Create(&tran).Error
	return
}
