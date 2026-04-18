package external

import (
	"context"
	"errors"

	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	pb "github.com/mhasnanr/ewallet-wallet/cmd/tokenvalidation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGRPC struct{}

func (e *UserGRPC) ValidateToken(ctx context.Context, accessToken string) (*pb.TokenResponse, error) {
	var (
		serverAddr = bootstrap.GetEnv("USER_GRPC_URL", "")
	)

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("failed to dial ums grpc")
	}

	defer conn.Close()

	client := pb.NewTokenValidationClient(conn)
	req := &pb.TokenRequest{
		Token: accessToken,
	}

	response, err := client.ValidateToken(ctx, req)
	if err != nil {
		return nil, errors.New("failed to validate token")
	}

	return response, nil
}
