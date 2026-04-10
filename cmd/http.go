package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	"github.com/mhasnanr/ewallet-wallet/internal/handler"
	"github.com/mhasnanr/ewallet-wallet/internal/repository"
	"github.com/mhasnanr/ewallet-wallet/internal/services"
	"gorm.io/gorm"
)

func ServeHTTP(db *gorm.DB) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "server is healthy"})
	})

	walletRepository := repository.NewWalletRepository(db)
	walletService := services.NeWalletService(walletRepository)
	walletHandler := handler.NewWalletHandler(walletService)

	walletHandler.RegisterRoute(r)

	server := &http.Server{Addr: ":" + bootstrap.GetEnv("HTTP_PORT", "8080"), Handler: r}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server stopped")
	}
}
