package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gitlab.com/investio/backend/sim-api/v1/model"
	"gitlab.com/investio/backend/sim-api/v1/service"
)

type TransactionController interface {
	GetTransaction(ctx *gin.Context)
}

type transactionController struct {
	authService        service.AuthService
	transactionService service.TransactionService
}

func NewTransactionController(auth service.AuthService, transaction service.TransactionService) TransactionController {
	return &transactionController{
		authService:        auth,
		transactionService: transaction,
	}
}

func (c *transactionController) GetTransaction(ctx *gin.Context) {
	var (
		transList []model.Transaction
	)

	// Get access token
	accessJWT, errReason := c.authService.ValidateAccessToken(ctx.Request)
	if errReason != "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"reason": errReason,
		})
		return
	}

	log.Info(transList)

	if err := c.transactionService.Get(&transList, accessJWT.UserID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"reason": "Unable to get transaction",
		})
		return
	}

	// log.Info(transList)

	ctx.JSON(http.StatusOK, transList)
}
