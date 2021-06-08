package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitlab.com/investio/backend/sim-api/db"
	"gitlab.com/investio/backend/sim-api/v1/controller"
	"gitlab.com/investio/backend/sim-api/v1/service"
)

var (
	log = logrus.New()

	authService        = service.NewAuthService()
	portService        = service.NewPortService()
	walletService      = service.NewWalletService()
	transactionService = service.NewTransctionService()

	portController        = controller.NewPortController(authService, portService, walletService, transactionService)
	walletController      = controller.NewWalletController(authService, walletService)
	transactionController = controller.NewTransactionController(authService, transactionService)
)

func getVersion(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": "1.0",
	})
}

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Warn("Main: Not using .env file")
		}
	}

	if err := db.SetupDB(); err != nil {
		log.Panic(err)
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:2564", "http://192.168.50.121:3003", "https://investio.dewkul.me", "https://investio.netlify.app"}
	// corsConfig.AllowMethods = []string{"PUT"}
	corsConfig.AllowHeaders = []string{"Authorization", "content-type"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for VueJS
	corsConfig.AddAllowMethods("OPTIONS")
	r.Use(cors.New(corsConfig))

	v1 := r.Group("/sim/v1")
	{
		v1.GET("/port", portController.GetFundsInPort)
		p := v1.Group("/port")
		{
			p.POST("/buy", portController.BuyFund)
			p.POST("/sell", portController.SellFund)
		}
		v1.GET("/wallet", walletController.GetWallet)
		v1.GET("/orders", transactionController.GetTransaction)
		v1.GET("/ver", getVersion)
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5005"
	}
	log.Panic(r.Run(":" + port))
}
