package constants

import "net/http"

var (
	ErrBadRequest   = "bad request"
	ErrUnauthorized = "unauthorized"
)

var (
	ErrUserIDRequired               = "user id is required"
	ErrReferenceRequired            = "reference is required"
	ErrAmountRequired               = "amount is required"
	ErrAmountMustBeGreaterThanZero  = "amount must be greater than zero"
	ErrUserNotFound                 = "user not found"
	ErrFailedToGetBalance           = "failed to get balance"
	ErrFailedToCreateWallet         = "failed to create wallet"
	ErrFailedToGetUserData          = "failed to get user data"
	ErrFailedToParseToken           = "failed to parse token"
	ErrFailedToParseUser            = "failed to parse user"
	ErrFailedToGetWalletTransaction = "failed to get wallet transaction"
	ErrDuplicateReference           = "reference is duplicate"
	ErrFailedToUpdateBalance        = "failed to update balance"
	ErrFailedToInsertTransaction    = "failed to insert transaction"
)

var ValidationErrorMap = map[string]map[string]string{
	"required": {},
}

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

var (
	ErrorUnauthorized                 = NewAppError(http.StatusUnauthorized, ErrUnauthorized)
	ErrorBadRequest                   = NewAppError(http.StatusBadRequest, ErrBadRequest)
	ErrorUserIDRequired               = NewAppError(http.StatusBadRequest, ErrUserIDRequired)
	ErrorAmountRequired               = NewAppError(http.StatusBadRequest, ErrAmountRequired)
	ErrorAmountMustBeGreaterThanZero  = NewAppError(http.StatusBadRequest, ErrAmountMustBeGreaterThanZero)
	ErrorUserNotFound                 = NewAppError(http.StatusNotFound, ErrUserNotFound)
	ErrorFailedToGetBalance           = NewAppError(http.StatusInternalServerError, ErrFailedToGetBalance)
	ErrorFailedToCreateWallet         = NewAppError(http.StatusInternalServerError, ErrFailedToCreateWallet)
	ErrorFailedToGetUserData          = NewAppError(http.StatusUnauthorized, ErrFailedToGetUserData)
	ErrorFailedToParseToken           = NewAppError(http.StatusUnauthorized, ErrFailedToParseToken)
	ErrorFailedToParseUser            = NewAppError(http.StatusInternalServerError, ErrFailedToParseUser)
	ErrorFailedToGetWalletTransaction = NewAppError(http.StatusNotFound, ErrFailedToGetWalletTransaction)
	ErrorDuplicateReference           = NewAppError(http.StatusConflict, ErrDuplicateReference)
	ErrorFailedToUpdateBalance        = NewAppError(http.StatusInternalServerError, ErrFailedToUpdateBalance)
	ErrorFailedToInsertTransaction    = NewAppError(http.StatusInternalServerError, ErrFailedToInsertTransaction)
)
