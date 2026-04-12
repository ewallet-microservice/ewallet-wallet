package services

import (
	"context"

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

type WalletService struct {
	repo    WalletRepository
	userAPI UserAPI
}

func NeWalletService(repo WalletRepository, userAPI UserAPI) *WalletService {
	return &WalletService{repo, userAPI}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	return s.repo.CreateWallet(ctx, wallet)
}

func (s *WalletService) GetBalance(ctx context.Context, accessToken string) (models.BalanceResponse, error) {
	var response models.BalanceResponse

	user, err := s.userAPI.ValidateToken(accessToken)
	if err != nil {
		return response, err
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, user.UserID)
	if err != nil {
		return response, err
	}

	response.Balance = wallet.Balance

	return response, nil
}
