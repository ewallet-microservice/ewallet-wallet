package grpc

import (
	"context"

	pb "github.com/mhasnanr/ewallet-wallet/cmd/walletTransaction"
	"github.com/mhasnanr/ewallet-wallet/constants"
	"github.com/mhasnanr/ewallet-wallet/internal/helpers"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletTransactionService interface {
	DebitBalance(ctx context.Context, userID int, req models.TransactionRequest) (models.BalanceResponse, error)
	CreditBalance(ctx context.Context, userID int, req models.TransactionRequest) (models.BalanceResponse, error)
}

type WalletTransactionSystem struct {
	pb.UnimplementedWalletTransactionServer
	svc WalletTransactionService
}

func NewWalletTransactionSystem(walletSvc WalletTransactionService) *WalletTransactionSystem {
	return &WalletTransactionSystem{svc: walletSvc}
}

func (e *WalletTransactionSystem) DebitBalance(ctx context.Context, request *pb.WalletRequest) (*pb.WalletResponse, error) {
	userID := request.GetUserId()
	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	walletRequest := models.TransactionRequest{
		Reference: request.GetReference(),
		Amount:    request.GetAmount(),
	}
	wallet, err := e.svc.DebitBalance(ctx, int(request.GetUserId()), walletRequest)
	if err != nil {
		return nil, helpers.MapAppErrorToGRPC(err)
	}

	return &pb.WalletResponse{
		Message: constants.DebitBalance,
		Data: &pb.WalletData{
			Balance: wallet.Balance,
		},
	}, nil
}

func (e *WalletTransactionSystem) CreditBalance(ctx context.Context, request *pb.WalletRequest) (*pb.WalletResponse, error) {
	userID := request.GetUserId()
	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	walletRequest := models.TransactionRequest{
		Reference: request.GetReference(),
		Amount:    request.GetAmount(),
	}
	wallet, err := e.svc.CreditBalance(ctx, int(request.GetUserId()), walletRequest)
	if err != nil {
		return nil, helpers.MapAppErrorToGRPC(err)
	}

	return &pb.WalletResponse{
		Message: constants.CreditBalance,
		Data: &pb.WalletData{
			Balance: wallet.Balance,
		},
	}, nil
}
