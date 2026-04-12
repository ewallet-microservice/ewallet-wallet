package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mhasnanr/ewallet-wallet/constants"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func SendResponseHTTP(c *gin.Context, code int, message string, data any) {
	resp := Response{
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}

func ConstructErrString(errors validator.ValidationErrors) string {
	errStrings := make([]string, len(errors))

	for i := range errors {
		var error = errors[i]
		var errMsg = constants.ValidationErrorMap[error.Tag()][error.Namespace()]
		errStrings[i] = errMsg
	}

	return strings.Join(errStrings, ", ")
}
