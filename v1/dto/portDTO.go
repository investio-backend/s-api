package dto

import (
	"github.com/shopspring/decimal"
	"gitlab.com/investio/backend/sim-api/v1/model"
)

type OrderRequest struct {
	DataDate model.Date      `json:"date"`
	PortID   uint            `json:"port_id"`
	FundID   string          `json:"fund_id"`
	FundCode string          `json:"fund_code"`
	BcatID   uint8           `json:"bcat_id"`
	Amount   decimal.Decimal `json:"amount"`
	Unit     decimal.Decimal `json:"unit"`
	NAV      decimal.Decimal `json:"nav"`
}
