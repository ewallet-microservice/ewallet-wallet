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
	CreditBalance(ctx context.Context, userID int, request models.TransactionRequest) (models.BalanceResponse, error)
}

type AuthMiddleware interface {
	MiddlewareAccessToken(c *gin.Context)
}

type WalletHandler struct {
	svc            WalletService
	authMiddleware AuthMiddleware
}

func NewWalletHandler(svc WalletService, authMiddleware AuthMiddleware) *WalletHandler {
	return &WalletHandler{svc, authMiddleware}
}

func (h *WalletHandler) RegisterRoute(r *gin.Engine) {
	walletV1 := r.Group("/wallets/v1")
	walletV1.POST("/", h.createWallet)
	walletV1.GET("/balance", h.authMiddleware.MiddlewareAccessToken, h.getBalance)
	walletV1.PUT("/balance/credit", h.authMiddleware.MiddlewareAccessToken, h.creditBalance)
}

func (h *WalletHandler) createWallet(c *gin.Context) {
	var req models.Wallet

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.writeErrorResponse(c, constants.ErrorBadRequest, nil)
		return
	}

	if req.UserID == 0 {
		h.writeErrorResponse(c, constants.ErrorUserIDRequired, nil)
		return
	}

	if err = req.Validate(); err != nil {
		h.writeErrorResponse(c, err, nil)
		return
	}

	err = h.svc.CreateWallet(c.Request.Context(), &req)
	if err != nil {
		h.writeErrorResponse(c, constants.ErrorFailedToCreateWallet, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.WalletCreated, req)
}

func (h *WalletHandler) getBalance(c *gin.Context) {
	userData, ok := c.Get("tokenData")

	if !ok {
		h.writeErrorResponse(c, constants.ErrorFailedToGetUserData, nil)
		return
	}

	data, ok := userData.(models.TokenData)
	if !ok {
		h.writeErrorResponse(c, constants.ErrorFailedToParseToken, nil)
		return
	}

	wallet, err := h.svc.GetBalance(c.Request.Context(), data.UserID)
	if err != nil {
		h.writeErrorResponse(c, constants.ErrorFailedToGetBalance, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.GetBalance, wallet)
}

func (h *WalletHandler) creditBalance(c *gin.Context) {
	var req models.TransactionRequest

	userData, ok := c.Get("tokenData")

	if !ok {
		h.writeErrorResponse(c, constants.ErrorFailedToGetUserData, nil)
		return
	}

	data, ok := userData.(models.TokenData)
	if !ok {
		h.writeErrorResponse(c, constants.ErrorFailedToParseToken, nil)
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.writeErrorResponse(c, constants.ErrorBadRequest, nil)
		return
	}

	if err = req.Validate(); err != nil {
		h.writeErrorResponse(c, err, nil)
		return
	}

	resp, err := h.svc.CreditBalance(c.Request.Context(), data.UserID, req)
	if err != nil {
		h.writeErrorResponse(c, err, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.CreditBalance, resp)
}

func (h *WalletHandler) writeErrorResponse(c *gin.Context, err error, data any) {
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
