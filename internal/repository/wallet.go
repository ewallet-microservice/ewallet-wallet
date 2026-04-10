package repository

import (
	"context"

	"github.com/mhasnanr/ewallet-wallet/internal/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	err := r.db.Create(wallet).Error
	if err != nil {
		return err
	}

	return nil
}
