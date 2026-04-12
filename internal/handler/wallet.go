package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mhasnanr/ewallet-wallet/constants"
	"github.com/mhasnanr/ewallet-wallet/internal/helpers"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	GetBalance(ctx context.Context, userID int) (models.BalanceResponse, error)
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
	walletV1.GET("/balance", w.getBalance)
}

func (w *WalletHandler) createWallet(c *gin.Context) {
	var req models.Wallet

	err := c.ShouldBindJSON(&req)
	if err != nil {
		w.writeErrorResponse(c, constants.ErrorBadRequest, nil)
		return
	}

	if req.UserID == 0 {
		w.writeErrorResponse(c, constants.ErrorUserIDRequired, nil)
		return
	}

	if err = req.Validate(); err != nil {
		w.writeErrorResponse(c, err, nil)
		return
	}

	err = w.svc.CreateWallet(c.Request.Context(), &req)
	if err != nil {
		w.writeErrorResponse(c, constants.ErrorFailedToCreateWallet, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.WalletCreated, req)
}

func (w *WalletHandler) getBalance(c *gin.Context) {
	userID := c.GetInt("userID")

	if userID == 0 {
		w.writeErrorResponse(c, constants.ErrorUserIDRequired, nil)
		return
	}

	wallet, err := w.svc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		w.writeErrorResponse(c, constants.ErrorFailedToGetBalance, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.GetBalance, wallet)
}

func (w *WalletHandler) writeErrorResponse(c *gin.Context, err error, data any) {
	var appErr *constants.AppError
	var valErrs validator.ValidationErrors

	if errors.As(err, &appErr) {
		helpers.SendResponseHTTP(c, appErr.StatusCode, appErr.Message, data)
		return
	}

	if errors.As(err, &valErrs) {
		errStr := helpers.ConstructErrString(valErrs)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, errStr, data)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
}
