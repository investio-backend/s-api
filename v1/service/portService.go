package service

import (
	"errors"

	"github.com/shopspring/decimal"
	"gitlab.com/investio/backend/sim-api/db"
	"gitlab.com/investio/backend/sim-api/v1/dto"
	"gitlab.com/investio/backend/sim-api/v1/model"
)

type PortService interface {
	CreatePort(userID uint) (port model.Port, err error)
	GetPort(p *model.Port, userID uint) (err error)
	GetFunds(funds *[]model.PortFund, portID uint) (err error)
	AddOrUpdateFund(req dto.OrderRequest) (err error)
	RedeemFund(req dto.OrderRequest) (err error)
}

type portService struct {
}

func NewPortService() PortService {
	return &portService{}
}

func (s *portService) CreatePort(userID uint) (port model.Port, err error) {
	port = model.Port{
		PortName:           "My first port",
		UserID:             userID,
		ProfitLossRealized: decimal.NewFromInt(0),
		AllCost:            decimal.NewFromInt(0),
	}
	err = db.SimDB.Create(&port).Error
	return
}

func (s *portService) GetPort(p *model.Port, userID uint) (err error) {
	err = db.SimDB.Where("user_id = ?", userID).First(p).Error
	return

}

func (s *portService) GetFunds(funds *[]model.PortFund, portID uint) (err error) {
	err = db.SimDB.Where("port_id = ?", portID).Find(funds).Error
	return
}

func (s *portService) AddOrUpdateFund(req dto.OrderRequest) (err error) {
	var (
		fund model.PortFund
		port model.Port
	)
	if err = db.SimDB.Where("ID = ?", req.PortID).First(&port).Error; err != nil {
		if (model.Port{}) == port {
			return errors.New("port not found")
		}
		return
	}

	port.AllCost = port.AllCost.Add(req.Amount)
	if err = db.SimDB.Save(&port).Error; err != nil {
		return
	}

	if err = db.SimDB.Where("fund_code = ?", req.FundCode).Where("port_id = ?", req.PortID).First(&fund).Error; err != nil {
		// Create
		fund := model.PortFund{
			FundID:   req.FundID,
			FundCode: req.FundCode,
			BcatID:   req.BcatID,
			PortID:   req.PortID,
			Cost:     req.Amount,
			Unit:     req.Unit,
		}
		err = db.SimDB.Create(&fund).Error
		return
	}

	// Update
	fund.Cost = fund.Cost.Add(req.Amount)
	fund.Unit = fund.Unit.Add(req.Unit)
	fund.BcatID = req.BcatID
	err = db.SimDB.Save(&fund).Error
	return
}

func (s *portService) RedeemFund(req dto.OrderRequest) (err error) {
	var (
		fund model.PortFund
		port model.Port
	)
	if err = db.SimDB.First(&port, req.PortID).Error; err != nil {
		if (model.Port{}) == port {
			return errors.New("port not found")
		}
		return
	}

	port.AllCost = port.AllCost.Add(req.Amount)
	if err = db.SimDB.Save(&port).Error; err != nil {
		return
	}

	if err = db.SimDB.Where("fund_code = ?", req.FundCode).Where("port_id = ?", req.PortID).First(&fund).Error; err != nil {
		return
	}

	if fund.Cost.Sub(req.Amount).LessThan(decimal.NewFromInt(0)) {
		return errors.New("reject: amount will be less than 0")
	}

	if fund.Unit.Sub(req.Unit).LessThan(decimal.NewFromInt(0)) {
		return errors.New("reject: unit will be less than 0")
	}

	fund.Cost = fund.Cost.Sub(req.Amount)
	fund.Unit = fund.Unit.Sub(req.Unit)
	fund.BcatID = req.BcatID
	err = db.SimDB.Save(&fund).Error
	return
}
