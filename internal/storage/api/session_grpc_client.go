package api

import (
	"cbe-super-app-member-auth/internal/constants/errors"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	session "cbe-super-app-member-auth/grpc"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SessionGRPCClient struct {
	client session.SessionServiceClient
	conn   *grpc.ClientConn
	logger utils.Logger
}

func NewSessionGRPCClient(serverAddress string, logger utils.Logger) (session.SessionServiceClient, *SessionGRPCClient, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Failed to connect to gRPC server: %v", err)
		return nil, nil, fmt.Errorf("GRPC_CONNECTION_FAILED")
	}

	client := session.NewSessionServiceClient(conn)

	return client, &SessionGRPCClient{
		client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

// Close closes the gRPC connection
func (c *SessionGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateSession creates a new session via gRPC
func (c *SessionGRPCClient) CreateSession(ctx context.Context, request *session.CreateSessionRequest, opts ...grpc.CallOption) (*session.CreateSessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{
		"client": "member-auth",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := c.client.CreateSession(ctx, request)
	if err != nil {
		c.logger.Errorf("Failed to create session via gRPC: %v", err)
		return nil, errors.ErrSessionCreationFailed
	}

	c.logger.Infof("Successfully created session with ID: %s", response.GetId())
	return response, nil
}

// UpdateSession updates an existing session via gRPC
func (c *SessionGRPCClient) UpdateSession(ctx context.Context, request *session.UpdateSessionRequest, opts ...grpc.CallOption) (*session.UpdateSessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{
		"client": "member-auth",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := c.client.UpdateSession(ctx, request)
	if err != nil {
		c.logger.Errorf("Failed to update session via gRPC: %v", err)
		return nil, errors.ErrSessionCreationFailed
	}

	c.logger.Infof("Successfully updated session with ID: %s", response.GetSessionId())
	return response, nil
}

// GetSession retrieves session information via gRPC
func (c *SessionGRPCClient) GetSession(ctx context.Context, request *session.GetSessionRequest, opts ...grpc.CallOption) (*session.GetSessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{
		"client": "member-auth",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := c.client.GetSession(ctx, request)
	if err != nil {
		c.logger.Errorf("Failed to get session via gRPC: %v", err)
		return nil, errors.ErrSessionRetrivalFailed
	}

	c.logger.Infof("Successfully retrieved session with ID: %s", request.GetSessionId())
	return response, nil
}

// HealthCheck performs a health check via gRPC
func (c *SessionGRPCClient) HealthCheck(ctx context.Context, opts ...grpc.CallOption) (*session.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{
		"client": "member-auth",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	request := &session.HealthCheckRequest{}
	response, err := c.client.HealthCheck(ctx, request)
	if err != nil {
		c.logger.Errorf("Failed to perform health check via gRPC: %v", err)
		return nil, errors.ErrHealthCheckFailed
	}

	c.logger.Infof("Health check successful: %s", response.GetStatus())
	return response, nil
}
