package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	"github.com/mhasnanr/ewallet-wallet/external"
	"github.com/mhasnanr/ewallet-wallet/internal/handler"
	"github.com/mhasnanr/ewallet-wallet/internal/middleware"
	"github.com/mhasnanr/ewallet-wallet/internal/repository"
	"github.com/mhasnanr/ewallet-wallet/internal/services"
	"github.com/mhasnanr/ewallet-wallet/internal/transactor"
	"gorm.io/gorm"
)

func ServeHTTP(db *gorm.DB) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "server is healthy"})
	})

	userGRPCClient, grpcConn, err := external.NewUserGRPC()
	if err != nil {
		log.Fatalf("failed to initialized gRPC: %v", err)
	}

	defer grpcConn.Close()

	authMiddleware := middleware.NewAuthMiddleware(userGRPCClient)
	txManager := transactor.NewTransactor(db)
	walletRepository := repository.NewWalletRepository(db)
	walletService := services.NeWalletService(walletRepository, txManager)
	walletHandler := handler.NewWalletHandler(walletService, authMiddleware)

	walletHandler.RegisterRoute(r)

	httpPort := bootstrap.GetEnv("HTTP_PORT", "8081")
	server := &http.Server{Addr: ":" + httpPort, Handler: r}

	fmt.Printf("http server is running on port %s...\n", httpPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server stopped")
	}
}
