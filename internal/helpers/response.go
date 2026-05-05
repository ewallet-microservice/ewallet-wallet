package helpers

import (
	"fmt"
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
		var err = errors[i]
		if tagMap, ok := constants.ValidationErrorMap[err.Tag()]; ok {
			if msg, ok := tagMap[err.Namespace()]; ok && msg != "" {
				errStrings[i] = msg
				continue
			}
		}
		errStrings[i] = fmt.Sprintf("Field %s failed on %s validation", err.Field(), err.Tag())
	}

	return strings.Join(errStrings, ", ")
}
