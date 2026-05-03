package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type WalletTransaction struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	WalletID        int       `json:"wallet_id" gorm:"column:wallet_id"`
	Amount          int64     `json:"amount" gorm:"column:amount;type:bigint"`
	Reference       string    `json:"reference" gorm:"column:reference;unique"`
	TransactionType string    `json:"transaction_type" gorm:"column:transaction_type"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (*WalletTransaction) TableName() string {
	return "wallet_transactions"
}

type TransactionRequest struct {
	Reference string `json:"reference" gorm:"column:reference;unique" validate:"required"`
	Amount    int64  `json:"amount" gorm:"column:amount;type:bigint" validate:"required,gt=0"`
}

func (f *TransactionRequest) Validate() error {
	v := validator.New()
	return v.Struct(f)
}

type TransactionHistory struct {
	ID        int    `db:"id"`
	Reference string `db:"string"`
}
