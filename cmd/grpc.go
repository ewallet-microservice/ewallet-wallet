package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	"github.com/mhasnanr/ewallet-wallet/cmd/wallet"
	"github.com/mhasnanr/ewallet-wallet/cmd/walletTransaction"
	handler "github.com/mhasnanr/ewallet-wallet/internal/handler/grpc"
	"github.com/mhasnanr/ewallet-wallet/internal/repository"
	"github.com/mhasnanr/ewallet-wallet/internal/services"
	"github.com/mhasnanr/ewallet-wallet/internal/transactor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func ServeGRPC(db *gorm.DB) {
	grpcPort := bootstrap.GetEnv("GRPC_PORT", "7001")
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("failed to listen grpc port: ", err)
	}

	server := grpc.NewServer()
	txManager := transactor.NewTransactor(db)
	walletRepository := repository.NewWalletRepository(db)
	walletService := services.NewWalletService(walletRepository, txManager)
	walletSystem := handler.NewWalletSystem(walletService)
	walletTxSystem := handler.NewWalletTransactionSystem(walletService)

	wallet.RegisterWalletServer(server, walletSystem)
	walletTransaction.RegisterWalletTransactionServer(server, walletTxSystem)

	reflection.Register(server)

	fmt.Printf("gRPC server is running on port %s...\n", grpcPort)
	if err := server.Serve(listener); err != nil {
		log.Fatal("failed to serve grpc port: ", err)
	}
}
