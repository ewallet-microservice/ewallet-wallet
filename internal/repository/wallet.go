package repository

import (
	"context"
	"errors"

	"github.com/mhasnanr/ewallet-wallet/constants"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
	"github.com/mhasnanr/ewallet-wallet/internal/transactor"
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

func (r *WalletRepository) GetWalletTransactionByReference(ctx context.Context, reference string) (models.WalletTransaction, error) {
	var transaction models.WalletTransaction

	err := r.db.Where("reference = ?", reference).First(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *WalletRepository) GetWalletForLock(ctx context.Context, userID int) (models.Wallet, error) {
	var wallet models.Wallet

	err := r.getExecutor(ctx).Raw("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE", userID).Scan(&wallet).Error
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, userID int, amount int64) error {
	err := r.getExecutor(ctx).Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", amount, userID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) CreateWalletTransaction(ctx context.Context, walletTransaction *models.WalletTransaction) error {
	err := r.getExecutor(ctx).Create(walletTransaction).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) getExecutor(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactor.TxKey{}).(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}
