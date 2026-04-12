package constants

import "net/http"

var ErrBadRequest = "bad request"

var (
	ErrUserIDRequired       = "user id is required"
	ErrUserNotFound         = "user not found"
	ErrFailedToGetBalance   = "failed to get balance"
	ErrFailedToCreateWallet = "failed to create wallet"
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
	ErrorBadRequest           = NewAppError(http.StatusBadRequest, ErrBadRequest)
	ErrorUserIDRequired       = NewAppError(http.StatusBadRequest, ErrUserIDRequired)
	ErrorUserNotFound         = NewAppError(http.StatusNotFound, ErrUserNotFound)
	ErrorFailedToGetBalance   = NewAppError(http.StatusInternalServerError, ErrFailedToGetBalance)
	ErrorFailedToCreateWallet = NewAppError(http.StatusInternalServerError, ErrFailedToCreateWallet)
)
