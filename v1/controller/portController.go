package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gitlab.com/investio/backend/sim-api/v1/dto"
	"gitlab.com/investio/backend/sim-api/v1/model"
	"gitlab.com/investio/backend/sim-api/v1/service"
)

type PortController interface {
	GetFundsInPort(ctx *gin.Context)
	BuyFund(ctx *gin.Context)
	SellFund(ctx *gin.Context)
}

type portController struct {
	authService        service.AuthService
	portService        service.PortService
	walletService      service.WalletService
	transactionService service.TransactionService
}

func NewPortController(auth service.AuthService, port service.PortService, wallet service.WalletService, transaction service.TransactionService) PortController {
	return &portController{
		authService:        auth,
		portService:        port,
		walletService:      wallet,
		transactionService: transaction,
	}
}

func (c *portController) GetFundsInPort(ctx *gin.Context) {
	var (
		port        model.Port
		fundsInPort []model.PortFund
	)

	// Get access token
	accessJWT, errReason := c.authService.ValidateAccessToken(ctx.Request)
	if errReason != "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"reason": errReason,
		})
		return
	}

	if err := c.portService.GetPort(&port, accessJWT.UserID); err != nil {
		port, err = c.portService.CreatePort(accessJWT.UserID)
		if err != nil {
			log.Error("CREATE PORT IN GetPort ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Unable to create port",
			})
			return
		}
	}

	if err := c.portService.GetFunds(&fundsInPort, port.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"reason": "Unable to get funds in port",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"port_id":     port.ID,
		"port_name":   port.PortName,
		"pl_realized": port.ProfitLossRealized,
		"sum_cost":    port.AllCost,
		"funds":       fundsInPort,
	})
}

func (c *portController) BuyFund(ctx *gin.Context) {
	var (
		req  dto.OrderRequest
		port model.Port
	)

	// Get access token
	accessJWT, errReason := c.authService.ValidateAccessToken(ctx.Request)
	if errReason != "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"reason": errReason,
		})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Invalid data provided",
		})
		return
	}

	if err := c.portService.GetPort(&port, accessJWT.UserID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Read port failed",
		})
		return
	}

	if port.ID != req.PortID {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid request",
		})
		return
	}

	if err := c.walletService.Purchase(req.Amount, accessJWT.UserID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Purchase failed: " + err.Error(),
		})
		return
	}

	if err := c.portService.AddOrUpdateFund(req); err != nil {
		if err := c.walletService.ReversePurchase(req.Amount, accessJWT.UserID); err != nil {
			log.Error("Critial [AddOrUpdateFund] - <rev> wallet purchase failed ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Critical - <rev> purchase failed: " + err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Add/update failed: " + err.Error(),
		})
		return
	}

	transaction := model.Transaction{
		DataDate: req.DataDate.ParseTime(),
		PortID:   req.PortID,
		FundID:   req.FundID,
		FundCode: req.FundCode,
		BcatID:   req.BcatID,
		Type:     1, // buy
		UserID:   accessJWT.UserID,
		NAV:      req.NAV,
		Amount:   req.Amount,
		Unit:     req.Unit,
	}
	if err := c.transactionService.Write(&transaction); err != nil {
		if err := c.walletService.ReversePurchase(req.Amount, accessJWT.UserID); err != nil {
			log.Error("Critial [AddOrUpdateFund] - <rev> wallet purchase failed ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Critical - <rev> purchase failed: " + err.Error(),
			})
			return
		}

		// Reverse purchase fund
		if err := c.portService.RedeemFund(req); err != nil {
			log.Error("Critial [AddOrUpdateFund] - <rev> port purchase fund failed ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Critical - <rev> purchase fund failed: " + err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Add transaction failed: " + err.Error(),
		})
		return
	}

	// log.Info(accessJWT, req.Amount.Mul(decimal.NewFromInt(5)))
	ctx.Status(200)
}

func (c *portController) SellFund(ctx *gin.Context) {
	var (
		req  dto.OrderRequest
		port model.Port
	)
	// Get access token
	accessJWT, errReason := c.authService.ValidateAccessToken(ctx.Request)
	if errReason != "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"reason": errReason,
		})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Invalid data provided",
		})
		return
	}

	// Validate input
	if req.Amount.LessThanOrEqual(decimal.NewFromInt32(0)) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.Unit.LessThanOrEqual(decimal.NewFromInt32(0)) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.portService.GetPort(&port, accessJWT.UserID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Read port failed",
		})
		return
	}

	if port.ID != req.PortID {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid request",
		})
		return
	}

	if err := c.walletService.Redeem(req.Amount, accessJWT.UserID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Redeem failed: " + err.Error(),
		})
		return
	}

	if err := c.portService.RedeemFund(req); err != nil {
		if err := c.walletService.ReverseRedeem(req.Amount, accessJWT.UserID); err != nil {
			log.Error("Critial [RedeemFund] - <rev> wallet redeem failed ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Critical - <rev> redeem failed: " + err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Redeem fund failed: " + err.Error(),
		})
		return
	}

	transaction := model.Transaction{
		DataDate: req.DataDate.ParseTime(),
		PortID:   req.PortID,
		FundID:   req.FundID,
		FundCode: req.FundCode,
		BcatID:   req.BcatID,
		Type:     2, // buy
		UserID:   accessJWT.UserID,
		NAV:      req.NAV,
		Amount:   req.Amount,
		Unit:     req.Unit,
	}
	if err := c.transactionService.Write(&transaction); err != nil {
		if err := c.walletService.ReverseRedeem(req.Amount, accessJWT.UserID); err != nil {
			log.Error("Critial [RedeemFund] - <rev> wallet redeem failed ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Critical - <rev> redeem failed: " + err.Error(),
			})
			return
		}
		// Reverse redeem fund
		if err := c.portService.AddOrUpdateFund(req); err != nil {
			log.Error("Critial [RedeemFund] - <rev> port add/update failed ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Critical - <rev> redeem fund failed: " + err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Add transaction failed: " + err.Error(),
		})
		return
	}

	ctx.Status(200)
}

// func (c *portController) PiePortUnit(ctx *gin.Context) {
// 	var (
// 		port        model.Port
// 		fundsInPort []model.PortFund
// 	)
// 	// Get access token
// 	accessJWT, errReason := c.authService.ValidateAccessToken(ctx.Request)
// 	if errReason != "" {
// 		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
// 			"reason": errReason,
// 		})
// 		return
// 	}

// 	if err := c.portService.GetPort(&port, accessJWT.UserID); err != nil {
// 		port, err = c.portService.CreatePort(accessJWT.UserID)
// 		if err != nil {
// 			log.Error("CREATE PORT IN PiePortUnit ", err.Error())
// 			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
// 				"reason": "Unable to create port",
// 			})
// 			return
// 		}
// 	}

// 	if err := c.portService.GetFunds(&fundsInPort, port.ID); err != nil {
// 		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
// 			"reason": "Unable to get funds in port",
// 		})
// 		return
// 	}

// }
