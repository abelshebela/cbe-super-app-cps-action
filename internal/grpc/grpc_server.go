package grpc

import (
	"context"
	"net"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// GRPCServer wraps the gRPC server and its services
type GRPCServer struct {
	server     *grpc.Server
	bankServer *server.BankServer
	logger     utils.Logger
}

// NewGRPCServer creates a new gRPC server with bank service
func NewGRPCServer(bankHandler bank.BankHandlerService, logger utils.Logger) *GRPCServer {
	// Create gRPC server with options
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Create bank server
	bankServer := server.NewBankServer(bankHandler, logger)

	// Register services (we'll implement the actual registration when we have protobuf generated code)
	// For now, we'll create a custom registration
	
	// Enable reflection for easier debugging
	reflection.Register(grpcServer)

	return &GRPCServer{
		server:     grpcServer,
		bankServer: bankServer,
		logger:     logger,
	}
}

// Start starts the gRPC server on the specified port
func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		s.logger.Errorf("Failed to listen on port %s: %v", port, err)
		return err
	}

	s.logger.Infof("🚀 gRPC Server starting on port %s", port)
	
	if err := s.server.Serve(lis); err != nil {
		s.logger.Errorf("Failed to serve gRPC server: %v", err)
		return err
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *GRPCServer) Stop() {
	s.logger.Infof("Stopping gRPC server...")
	s.server.GracefulStop()
	s.logger.Infof("gRPC server stopped")
}

// loggingInterceptor provides logging for gRPC calls
func loggingInterceptor(logger utils.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infof("gRPC call: %s", info.FullMethod)
		
		resp, err := handler(ctx, req)
		
		if err != nil {
			logger.Errorf("gRPC call %s failed: %v", info.FullMethod, err)
		} else {
			logger.Infof("gRPC call %s completed successfully", info.FullMethod)
		}
		
		return resp, err
	}
}
