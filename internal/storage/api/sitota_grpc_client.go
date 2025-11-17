package api

import (
	"context"
	"fmt"
	"time"

	transactionpb "cbe-super-app-cps-action/grpc/sitota"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cc     *grpc.ClientConn
	client transactionpb.TransactionServiceClient
}

func NewSitotagRPCClient(ctx context.Context, logger utils.Logger, gRPCAddress string) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cc, err := grpc.NewClient(
		gRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial CBE service at %s: %w", gRPCAddress, err)
	}

	transactionClient := transactionpb.NewTransactionServiceClient(cc)
	// cbeClient := NewCbeToCbeServiceClient(cc)

	logger.Infof("Created new gRPC client for CBE-to-Cbe service")

	return &Client{
		cc:     cc,
		client: transactionClient,
	}, nil
}

func (c *Client) Close() error {
	if c.cc != nil {
		return c.cc.Close()
	}
	return nil
}

func (c *Client) GetTransaction(ctx context.Context, in *transactionpb.TransactionRequest, opts ...grpc.CallOption) (*transactionpb.Transaction, error) {
	return c.client.GetTransaction(ctx, in, opts...)
}

func (c *Client) GetTransactionsList(ctx context.Context, in *transactionpb.ListTransactionRequest, opts ...grpc.CallOption) (*transactionpb.TransactionList, error) {
	return c.client.GetTransactionsList(ctx, in, opts...)
}

func (c *Client) HealthCheck(ctx context.Context, in *transactionpb.HealthCheckRequest, opts ...grpc.CallOption) (*transactionpb.HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, in, opts...)
}
