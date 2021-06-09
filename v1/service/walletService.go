package service

import (
	"errors"

	"github.com/shopspring/decimal"
	"gitlab.com/investio/backend/sim-api/db"
	"gitlab.com/investio/backend/sim-api/v1/model"
)

type WalletService interface {
	GetWallet(wallet *model.Wallet, userID uint) (err error)
	CreateWallet(userID uint) (wallet model.Wallet, err error)
	Purchase(amount decimal.Decimal, userID uint) (err error)
	Redeem(amount decimal.Decimal, userID uint) (err error)
	ReversePurchase(amount decimal.Decimal, userID uint) (err error)
	ReverseRedeem(amount decimal.Decimal, userID uint) (err error)
}

type walletService struct {
	startBalance decimal.Decimal
}

func NewWalletService() WalletService {
	return &walletService{
		startBalance: decimal.NewFromInt32(1000000),
	}
}

func (s *walletService) CreateWallet(userID uint) (wallet model.Wallet, err error) {
	wallet = model.Wallet{
		UserID:       userID,
		AvaliableBal: s.startBalance,
		InOrderBal:   decimal.NewFromInt32(0),
		InAssetBal:   decimal.NewFromInt32(0),
		TotalSpend:   decimal.NewFromInt32(0),
	}
	err = db.SimDB.Create(&wallet).Error
	return
}

func (s *walletService) GetWallet(wallet *model.Wallet, userID uint) (err error) {
	err = db.SimDB.Where("user_id = ?", userID).First(&wallet).Error
	return
}

func (s *walletService) Purchase(amount decimal.Decimal, userID uint) (err error) {
	var wallet model.Wallet

	if err = db.SimDB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return
	}

	if wallet.AvaliableBal.Sub(amount).LessThan(decimal.NewFromInt(0)) {
		return errors.New("reject: avaliable balance will be less than zero")
	}

	wallet.AvaliableBal = wallet.AvaliableBal.Sub(amount)
	wallet.InAssetBal = wallet.InAssetBal.Add(amount)
	wallet.TotalSpend = wallet.TotalSpend.Add(amount)
	err = db.SimDB.Save(&wallet).Error
	return
}

func (s *walletService) ReversePurchase(amount decimal.Decimal, userID uint) (err error) {
	var wallet model.Wallet

	if err = db.SimDB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return
	}

	wallet.AvaliableBal = wallet.AvaliableBal.Add(amount)
	wallet.InAssetBal = wallet.InAssetBal.Sub(amount)
	wallet.TotalSpend = wallet.TotalSpend.Sub(amount)
	err = db.SimDB.Save(&wallet).Error
	return
}

func (s *walletService) Redeem(amount decimal.Decimal, userID uint) (err error) {
	var wallet model.Wallet

	if err = db.SimDB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return
	}

	if wallet.InAssetBal.Sub(amount).LessThan(decimal.NewFromInt(0)) {
		return errors.New("reject: in asset balance will be less than zero")
	}

	wallet.AvaliableBal = wallet.AvaliableBal.Add(amount)
	wallet.InAssetBal = wallet.InAssetBal.Sub(amount)
	err = db.SimDB.Save(&wallet).Error
	return
}

func (s *walletService) ReverseRedeem(amount decimal.Decimal, userID uint) (err error) {
	var wallet model.Wallet

	if err = db.SimDB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return
	}

	wallet.AvaliableBal = wallet.AvaliableBal.Sub(amount)
	wallet.InAssetBal = wallet.InAssetBal.Add(amount)
	err = db.SimDB.Save(&wallet).Error
	return
}
