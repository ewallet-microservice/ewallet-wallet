package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	pb "github.com/mhasnanr/ewallet-wallet/cmd/tokenvalidation"
	"github.com/mhasnanr/ewallet-wallet/internal/helpers"
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
	if auth == "" {
		log.Infow("authorization empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	token := strings.Split(auth, "Bearer ")[1]
	if token == "" {
		log.Infow("invalid token")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	user, err := m.userGRPC.ValidateToken(c.Request.Context(), token)
	if err != nil {
		log.Infow("invalid token")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	userData := user.GetData()
	if userData == nil {
		log.Infow("failed to parse user data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to parse user", nil)
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
