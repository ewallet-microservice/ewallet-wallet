package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	"github.com/mhasnanr/ewallet-wallet/external"
	"github.com/mhasnanr/ewallet-wallet/internal/handler"
	"github.com/mhasnanr/ewallet-wallet/internal/middleware"
	"github.com/mhasnanr/ewallet-wallet/internal/repository"
	"github.com/mhasnanr/ewallet-wallet/internal/services"
	"gorm.io/gorm"
)

func ServeHTTP(db *gorm.DB) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "server is healthy"})
	})

	userAPIExternal := &external.ExternalUserAPI{}
	authMiddleware := middleware.NewAuthMiddleware()

	walletRepository := repository.NewWalletRepository(db)
	walletService := services.NeWalletService(walletRepository, userAPIExternal)
	walletHandler := handler.NewWalletHandler(walletService, authMiddleware)

	walletHandler.RegisterRoute(r)

	server := &http.Server{Addr: ":" + bootstrap.GetEnv("HTTP_PORT", "8081"), Handler: r}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server stopped")
	}
}
