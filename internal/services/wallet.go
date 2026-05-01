package services

import (
	"context"
	"errors"

	"github.com/mhasnanr/ewallet-wallet/constants"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	GetWalletByUserID(ctx context.Context, userID int) (models.Wallet, error)
	GetWalletTransactionByReference(ctx context.Context, reference string) (models.WalletTransaction, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) (models.Wallet, error)
	CreateWalletTransaction(ctx context.Context, walletTransaction *models.WalletTransaction) error
	GetWalletForLock(ctx context.Context, userID int) (models.Wallet, error)
}

type TxManager interface {
	WithinTransaction(ctx context.Context, txFunc func(context.Context) error) error
}

type WalletService struct {
	repo      WalletRepository
	txManager TxManager
}

func NewWalletService(repo WalletRepository, txManager TxManager) *WalletService {
	return &WalletService{repo, txManager}
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

func (s *WalletService) CreditBalance(ctx context.Context, userID int, req models.TransactionRequest) (models.BalanceResponse, error) {
	var response models.BalanceResponse

	history, err := s.repo.GetWalletTransactionByReference(ctx, req.Reference)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response, constants.ErrorFailedToGetWalletTransaction
		}
	}

	if history.ID > 0 {
		return response, constants.ErrorDuplicateReference
	}

	var wallet models.Wallet

	err = s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		wallet, err = s.repo.GetWalletForLock(ctx, userID)
		if err != nil {
			return err
		}

		wallet, err = s.repo.UpdateBalance(ctx, userID, req.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response, constants.ErrorFailedToUpdateBalance
	}

	walletTransaction := &models.WalletTransaction{
		WalletID:        wallet.ID,
		Amount:          req.Amount,
		Reference:       req.Reference,
		TransactionType: constants.CreditTransaction,
	}

	err = s.repo.CreateWalletTransaction(ctx, walletTransaction)
	if err != nil {
		return response, constants.ErrorFailedToInsertTransaction
	}

	response.Balance = wallet.Balance + req.Amount

	return response, nil
}

func (s *WalletService) DebitBalance(ctx context.Context, userID int, req models.TransactionRequest) (models.BalanceResponse, error) {
	var response models.BalanceResponse

	history, err := s.repo.GetWalletTransactionByReference(ctx, req.Reference)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response, constants.ErrorFailedToGetWalletTransaction
		}
	}

	if history.ID > 0 {
		return response, constants.ErrorDuplicateReference
	}

	var wallet models.Wallet

	err = s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		wallet, err = s.repo.GetWalletForLock(ctx, userID)
		if err != nil {
			return err
		}

		wallet, err = s.repo.UpdateBalance(ctx, userID, -req.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response, constants.ErrorFailedToUpdateBalance
	}

	walletTransaction := &models.WalletTransaction{
		WalletID:        wallet.ID,
		Amount:          req.Amount,
		Reference:       req.Reference,
		TransactionType: constants.DebitTransaction,
	}

	err = s.repo.CreateWalletTransaction(ctx, walletTransaction)
	if err != nil {
		return response, constants.ErrorFailedToInsertTransaction
	}

	response.Balance = wallet.Balance + req.Amount

	return response, nil
}
