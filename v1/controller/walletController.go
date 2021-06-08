package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/investio/backend/sim-api/v1/model"
	"gitlab.com/investio/backend/sim-api/v1/service"
)

type WalletController interface {
	GetWallet(ctx *gin.Context)
}

type walletController struct {
	authService   service.AuthService
	walletService service.WalletService
}

func NewWalletController(auth service.AuthService, wallet service.WalletService) WalletController {
	return &walletController{
		authService:   auth,
		walletService: wallet,
	}
}

func (c *walletController) GetWallet(ctx *gin.Context) {
	var (
		wallet model.Wallet
	)

	// Get access token
	accessJWT, errReason := c.authService.ValidateAccessToken(ctx.Request)
	if errReason != "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"reason": errReason,
		})
		return
	}

	if err := c.walletService.GetWallet(&wallet, accessJWT.UserID); err != nil {
		wallet, err = c.walletService.CreateWallet(accessJWT.UserID)
		if err != nil {
			log.Error("CREATE WALLET IN GetWallet ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"reason": "Unable to create wallet",
			})
			return
		}
	}

	ctx.JSON(200, wallet)
}
