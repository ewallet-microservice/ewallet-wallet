package services

import (
	"context"

	"github.com/mhasnanr/ewallet-wallet/internal/models"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	GetWalletByUserID(ctx context.Context, userID int) (models.Wallet, error)
}

type WalletService struct {
	repo WalletRepository
}

func NeWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	return s.repo.CreateWallet(ctx, wallet)
}

func (s *WalletService) GetBalance(ctx context.Context, userID int) (models.BalanceResponse, error) {
	var response models.BalanceResponse

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return response, err
	}

	response.Balance = wallet.Balance

	return response, nil
}
