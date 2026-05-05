package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	pb "github.com/mhasnanr/ewallet-wallet/cmd/tokenvalidation"
	"github.com/mhasnanr/ewallet-wallet/constants"
	"github.com/mhasnanr/ewallet-wallet/internal/models"
)

type UserServiceGRPC interface {
	ValidateToken(ctx context.Context, accessToken string) (*pb.TokenResponse, error)
}

type AuthMiddleware struct {
	userGRPC UserServiceGRPC
}

func NewAuthMiddleware(userGRPC UserServiceGRPC) *AuthMiddleware {
	return &AuthMiddleware{userGRPC}
}

func (m *AuthMiddleware) MiddlewareAccessToken(c *gin.Context) {
	var log = bootstrap.Log

	auth := c.Request.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		log.Infow("invalid or missing bearer token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		log.Infow("invalid token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	user, err := m.userGRPC.ValidateToken(c.Request.Context(), token)
	if err != nil {
		log.Infow("invalid token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	userData := user.GetData()
	if userData == nil {
		log.Infow("failed to parse user data")
		c.Error(constants.ErrorFailedToParseUser)
		c.Abort()
		return
	}

	tokenData := models.TokenData{
		UserID:   int(userData.UserId),
		Username: userData.Username,
		Fullname: userData.FullName,
		Email:    userData.Email,
	}

	c.Set("tokenData", tokenData)
	c.Next()
}
