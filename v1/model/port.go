package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Port struct {
	ID       uint   `gorm:"primaryKey" json:"port_id"`
	PortName string `json:"port_name"`
	UserID   uint   `json:"uid"`
	// ProfitLost        decimal.Decimal `sql:"type:decimal(14,4)"`
	// ProfitLostPercent decimal.Decimal `sql:"type:decimal(12,2)"`
	ProfitLossRealized decimal.Decimal `json:"pl_realized" gorm:"type:decimal(12,2);"`
	AllCost            decimal.Decimal `json:"sum_cost" gorm:"type:decimal(12,2);"`
	CreatedAt          time.Time       `json:"-"`
	UpdatedAt          time.Time       `json:"-"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName fund
func (Port) TableName() string {
	return "port"
}

type PortFund struct {
	ID         uint            `gorm:"primaryKey" json:"-"`
	FundID     string          `json:"fund_id"`
	FundCode   string          `json:"code"`
	BcatID     uint8           `json:"bcat_id"`
	PortID     uint            `json:"-"`
	Cost       decimal.Decimal `json:"cost" gorm:"type:decimal(12,2);"`
	Unit       decimal.Decimal `json:"unit" gorm:"type:decimal(14,4);"`
	PlRealized decimal.Decimal `json:"pl_realized" gorm:"type:decimal(12,2);"`
	CreatedAt  time.Time       `json:"-"`
	UpdatedAt  time.Time       `json:"-"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName fund
func (PortFund) TableName() string {
	return "port_fund"
}
