package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-wallet/internal/helpers"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
}

type WalletHandler struct {
	svc WalletService
}

func NewWalletHandler(svc WalletService) *WalletHandler {
	return &WalletHandler{svc}
}

func (w *WalletHandler) RegisterRoute(r *gin.Engine) {
	walletV1 := r.Group("/wallets/v1")
	walletV1.POST("/", w.createWallet)

}

func (w *WalletHandler) createWallet(c *gin.Context) {
	var req models.Wallet

	err := c.ShouldBindJSON(&req)
	if err != nil {
		helpers.SendResponseHTTP(c, http.StatusBadRequest, "bad request", nil)
		return
	}

	if req.UserID == 0 {
		helpers.SendResponseHTTP(c, http.StatusBadRequest, "user id is empty", nil)
		return
	}

	if err = req.Validate(); err != nil {
		fmt.Println(err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, "invalid body", nil)
		return
	}

	err = w.svc.CreateWallet(c.Request.Context(), &req)
	if err != nil {
		fmt.Println("dimari")
		fmt.Println(err)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, "wallet created", req)
}
