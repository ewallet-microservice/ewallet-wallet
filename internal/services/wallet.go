package services

import (
	"context"

	"github.com/mhasnanr/ewallet-wallet/internal/models"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
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
