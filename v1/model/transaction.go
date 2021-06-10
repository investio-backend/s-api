package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	ID        uint            `gorm:"primaryKey" json:"-"`
	DataDate  time.Time       `json:"data_date" gorm:"type:date;"`
	Type      uint32          `json:"transaction_type"` // 1-buy, 2-sell
	UserID    uint            `json:"-"`
	PortID    uint            `json:"port_id"`
	FundID    string          `json:"fund_id"`
	FundCode  string          `json:"code"`
	BcatID    uint8           `json:"bcat_id"`
	NAV       decimal.Decimal `gorm:"type:decimal(14,4);"`
	Amount    decimal.Decimal `json:"amount" gorm:"type:decimal(12,2);"`
	Unit      decimal.Decimal `json:"unit" gorm:"type:decimal(18,8);"`
	CreatedAt time.Time       `json:"timestamp"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName transaction
func (Transaction) TableName() string {
	return "transaction"
}
