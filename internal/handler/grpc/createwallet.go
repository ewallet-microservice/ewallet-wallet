package grpc

import (
	"context"

	pb "github.com/mhasnanr/ewallet-wallet/cmd/wallet"
	"github.com/mhasnanr/ewallet-wallet/constants"
	"github.com/mhasnanr/ewallet-wallet/internal/helpers"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
}

type WalletSystem struct {
	pb.UnimplementedWalletServer
	svc WalletService
}

func NewWalletSystem(walletSvc WalletService) *WalletSystem {
	return &WalletSystem{svc: walletSvc}
}

func (h *WalletSystem) CreateWallet(ctx context.Context, request *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	userID := request.GetUserId()
	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	var wallet = &models.Wallet{
		UserID: int(userID),
	}

	err := h.svc.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, helpers.MapAppErrorToGRPC(err)
	}

	return &pb.CreateWalletResponse{
		Message: constants.WalletCreated,
		Data: &pb.WalletData{
			UserId:  userID,
			Balance: wallet.Balance,
		},
	}, nil
}
