package service

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/grpc/server"
)

// BankServiceServer defines the gRPC service interface for bank operations
type BankServiceServer interface {
	GetBank(ctx context.Context, req *server.GetBankRequest) (*server.GetBankResponse, error)
	GetAllBanks(ctx context.Context, req *server.GetAllBanksRequest) (*server.GetAllBanksResponse, error)
	CreateBank(ctx context.Context, req *server.CreateBankRequest) (*server.CreateBankResponse, error)
	UpdateBank(ctx context.Context, req *server.UpdateBankRequest) (*server.UpdateBankResponse, error)
	DeleteBank(ctx context.Context, req *server.DeleteBankRequest) (*server.DeleteBankResponse, error)
	EnableDisableBank(ctx context.Context, req *server.EnableDisableBankRequest) (*server.EnableDisableBankResponse, error)
	UpdateBankLogo(ctx context.Context, req *server.UpdateBankLogoRequest) (*server.UpdateBankLogoResponse, error)
}

// Ensure BankServer implements BankServiceServer
var _ BankServiceServer = (*server.BankServer)(nil)
