package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Wallet struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"column:user_id;unique"`
	Balance   float64   `json:"balance" gorm:"column:balance;type:decimal(15,2)" validate:"gte=0"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (f *Wallet) Validate() error {
	v := validator.New()
	return v.Struct(f)
}
