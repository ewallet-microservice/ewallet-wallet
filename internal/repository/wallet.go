package repository

import (
	"context"
	"errors"

	"github.com/mhasnanr/ewallet-wallet/constants"
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

func (r *WalletRepository) GetWalletByUserID(ctx context.Context, userID int) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wallet, constants.ErrorUserNotFound
		}

		return wallet, err
	}

	return wallet, nil
}
