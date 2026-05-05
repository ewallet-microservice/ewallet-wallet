package helpers

import (
	"errors"
	"net/http"

	"github.com/mhasnanr/ewallet-wallet/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapAppErrorToGRPC(err error) error {
	var appErr *constants.AppError
	if errors.As(err, &appErr) {
		switch appErr.StatusCode {
		case http.StatusBadRequest:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case http.StatusUnauthorized:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case http.StatusForbidden:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case http.StatusNotFound:
			return status.Error(codes.NotFound, appErr.Message)
		case http.StatusConflict:
			return status.Error(codes.AlreadyExists, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	return status.Errorf(codes.Internal, "internal server error: %v", err)
}
