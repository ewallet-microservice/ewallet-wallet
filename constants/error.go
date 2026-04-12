package constants

import "net/http"

var (
	ErrBadRequest   = "bad request"
	ErrUnauthorized = "unauthorized"
)

var (
	ErrUserIDRequired       = "user id is required"
	ErrUserNotFound         = "user not found"
	ErrFailedToGetBalance   = "failed to get balance"
	ErrFailedToCreateWallet = "failed to create wallet"
	ErrFailedToGetToken     = "failed to get token"
	ErrFailedToParseToken   = "failed to parse token"
	ErrFailedToParseUser    = "failed to parse user"
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
	ErrorUnauthorized         = NewAppError(http.StatusUnauthorized, ErrUnauthorized)
	ErrorBadRequest           = NewAppError(http.StatusBadRequest, ErrBadRequest)
	ErrorUserIDRequired       = NewAppError(http.StatusBadRequest, ErrUserIDRequired)
	ErrorUserNotFound         = NewAppError(http.StatusNotFound, ErrUserNotFound)
	ErrorFailedToGetBalance   = NewAppError(http.StatusInternalServerError, ErrFailedToGetBalance)
	ErrorFailedToCreateWallet = NewAppError(http.StatusInternalServerError, ErrFailedToCreateWallet)
	ErrorFailedToGetToken     = NewAppError(http.StatusUnauthorized, ErrFailedToGetToken)
	ErrorFailedToParseToken   = NewAppError(http.StatusUnauthorized, ErrFailedToParseToken)
	ErrorFailedToParseUser    = NewAppError(http.StatusInternalServerError, ErrFailedToParseUser)
)
