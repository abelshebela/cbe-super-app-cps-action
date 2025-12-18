package client

import (
	cps_auth "cbe-super-app-cps-action/grpc/auth/proto"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGRPCClient struct {
	Client cps_auth.CpsAuthServiceClient
	conn   *grpc.ClientConn
	logger utils.Logger
}

func NewAuthGRPCClient(serverAddress string, logger utils.Logger) (*AuthGRPCClient, error) {
	address := strings.Split(serverAddress, ":")
	conn, err := grpc.Dial(address[0]+":"+address[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Failed to connect to gRPC server at:%s err:%v", serverAddress, err)
		return nil, fmt.Errorf("GRPC_CONNECTION_FAILED")
	}

	client := cps_auth.NewCpsAuthServiceClient(conn)

	return &AuthGRPCClient{
		Client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

// Close closes the gRPC connection
func (c *AuthGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
