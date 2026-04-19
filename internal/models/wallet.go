package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Wallet struct {
	ID        int       `json:"id" gorm:"primarykey"`
	UserID    int       `json:"user_id" gorm:"column:user_id;unique"`
	Balance   float64   `json:"balance" gorm:"column:balance;type:decimal(15,2)" validate:"gte=0"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (w *Wallet) Validate() error {
	v := validator.New()
	return v.Struct(w)
}

type WalletTransaction struct {
	ID              int       `json:"id" gorm:"primarykey"`
	WalletID        int       `json:"wallet_id" gorm:"column:wallet_id"`
	Amount          float64   `json:"amount" gorm:"column:amount;type:decimal(15, 2)"`
	Reference       string    `json:"reference" gorm:"column:reference;unique"`
	TransactionType string    `json:"transaction_type" gorm:"column:transaction_type;type:enum('CREDIT','DEBIT')"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (*WalletTransaction) TableName() string {
	return "wallet_transactions"
}
