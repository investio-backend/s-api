package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	ID           uint            `gorm:"primaryKey" json:"-"`
	AvaliableBal decimal.Decimal `json:"avalible_bal" gorm:"type:decimal(12,2);"`
	InOrderBal   decimal.Decimal `json:"inorder_bal" gorm:"type:decimal(12,2);"`
	InAssetBal   decimal.Decimal `json:"inasset_bal" gorm:"type:decimal(12,2);"`
	UserID       uint            `json:"-"`
	CreatedAt    time.Time       `json:"-"`
	UpdatedAt    time.Time       `json:"-"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName wallet
func (Wallet) TableName() string {
	return "wallet"
}
