package api

import (
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"google.golang.org/grpc"
)

type Client struct {
	cc *grpc.ClientConn
	// client
}

func NewSitotagRPCClient(ctx context.Context, logger utils.Logger, gRPCAddress string) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// cc, err := grpc.NewClient(
	// 	gRPCAddress,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithConnectParams(grpc.ConnectParams{
	// 		MinConnectTimeout: 5 * time.Second,
	// 	}),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to dial CBE service at %s: %w", gRPCAddress, err)
	// }

	// cbeClient := vault.NewCbeToCbeServiceClient(cc)

	logger.Infof("Created new gRPC client for CBE-to-Cbe service")

	// return &Client{
	// 	cc:     cc,
	// 	client: cbeClient,
	// }, nil

	return nil, nil
}

func (c *Client) Close() error {
	if c.cc != nil {
		return c.cc.Close()
	}
	return nil
}
