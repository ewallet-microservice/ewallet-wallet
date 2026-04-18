package services

import (
	"context"
	"errors"

	pb "github.com/mhasnanr/ewallet-wallet/cmd/tokenvalidation"
	"github.com/mhasnanr/ewallet-wallet/external"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	GetWalletByUserID(ctx context.Context, userID int) (models.Wallet, error)
}

type UserAPI interface {
	ValidateToken(string) (external.ValidateUserResponse, error)
}

type UserServiceGRPC interface {
	ValidateToken(ctx context.Context, accessToken string) (*pb.TokenResponse, error)
}

type WalletService struct {
	repo WalletRepository
	// userAPI UserAPI
	userSvcGRPC UserServiceGRPC
}

func NeWalletService(repo WalletRepository, userSvc UserServiceGRPC) *WalletService {
	return &WalletService{repo, userSvc}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	return s.repo.CreateWallet(ctx, wallet)
}

func (s *WalletService) GetBalance(ctx context.Context, accessToken string) (models.BalanceResponse, error) {
	var response models.BalanceResponse

	user, err := s.userSvcGRPC.ValidateToken(ctx, accessToken)
	if err != nil {
		return response, err
	}

	userData := user.GetData()
	if userData == nil {
		return response, errors.New("user data not found in response")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, int(userData.GetUserId()))
	if err != nil {
		return response, err
	}

	response.Balance = wallet.Balance

	return response, nil
}
