package server

import (
	bankpb "cbe-super-app-cps-action/grpc/bank"
	servicepb "cbe-super-app-cps-action/grpc/service/proto"
	walletpb "cbe-super-app-cps-action/grpc/wallet/proto"
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"log"
	"net"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"google.golang.org/grpc"
)

// server implements bankpb.BankServiceServer
// Embed UnimplementedBankServiceServer for forward compatibility

type server struct {
	bankpb.UnimplementedBankServiceServer
	walletpb.UnimplementedWalletServiceServer
	servicepb.UnimplementedServiceDetailsServiceServer
	bankHandler    service.BankService
	walletHandler  service.WalletService
	serviceHandler service.ServiceService
	logger         utils.Logger
}

func NewGrpcServer(bankHandler service.BankService, walletHandler service.WalletService, serviceHandler service.ServiceService, logger utils.Logger) *server {
	return &server{
		bankHandler:    bankHandler,
		walletHandler:  walletHandler,
		serviceHandler: serviceHandler,
		logger:         logger,
	}
}

// //////////////////////bank/////////////////
func (s *server) GetOneBank(ctx context.Context, req *bankpb.GetOneBankRequest) (*bankpb.GetOneBankResponse, error) {
	// TODO: Implement logic
	data, err := s.bankHandler.GetOneBank(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("Failed to get bank: %v", err)
		return nil, err
	}
	return &bankpb.GetOneBankResponse{Bank: s.bankMapper(data)}, nil
}

func (s *server) walletMapper(data *model.Wallet) *walletpb.Wallet {
	return &walletpb.Wallet{
		Id:        data.ID.Hex(),
		Name:      data.Name,
		Avatar:    data.Avatar,
		Code:      data.Code,
		IsDeleted: data.IsDeleted,
		Enabled:   data.Enabled,
	}
}

func (s *server) bankListMapper(data []*model.Bank) []*bankpb.Bank {
	var banks []*bankpb.Bank
	for _, bank := range data {
		banks = append(banks, s.bankMapper(bank))
	}
	return banks
}

func (s *server) GetAllBank(ctx context.Context, req *bankpb.GetAllBankRequest) (*bankpb.GetAllBankResponse, error) {

	data, err := s.bankHandler.GetAllBank(ctx, &types.Filter{Page: int(req.Page), PerPage: int(req.PerPage), Search: req.Search})
	if err != nil {
		s.logger.Errorf("Failed to get all banks: %v", err)
		return nil, err
	}

	return &bankpb.GetAllBankResponse{Banks: s.bankListMapper(data.Data), Metadata: buildPagination(data.Meta)}, nil
}
func (s *server) bankMapper(data *model.Bank) *bankpb.Bank {
	return &bankpb.Bank{
		Id:      data.ID.Hex(),
		Name:    data.Name,
		Bic:     data.BIC,
		Code:    data.Code,
		Logo:    data.Logo,
		Enabled: data.Enabled,
	}
}
func buildPagination(meta types.PaginationMeta) *bankpb.Meta {
	return &bankpb.Meta{
		TotalPages:  int32(meta.TotalPages),
		Limit:       int32(meta.Limit),
		TotalDocs:   int32(meta.TotalDocs),
		Page:        int32(meta.Page),
		HasNextPage: meta.HasNextPage,
		NextPage:    int32(*meta.NextPage),
		HasPrevPage: meta.HasPrevPage,
	}
}

// ///////////////////wallet///////////////////
func (s *server) GetAllWallet(ctx context.Context, req *walletpb.GetAllWalletRequest) (*walletpb.GetAllWalletResponse, error) {
	data, err := s.walletHandler.GetAllWallet(ctx, types.Filter{Page: int(req.Page), PerPage: int(req.PerPage), Search: req.Search})
	if err != nil {
		s.logger.Errorf("Failed to get all wallets: %v", err)
		return nil, err
	}
	return &walletpb.GetAllWalletResponse{Wallets: s.walletListMapper(data.Data), Metadata: buildPaginationWallet(data.Meta)}, nil
}

func (s *server) GetWallet(ctx context.Context, req *walletpb.GetWalletRequest) (*walletpb.GetWalletResponse, error) {
	data, err := s.walletHandler.GetWallet(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("Failed to get wallet: %v", err)
		return nil, err
	}
	return &walletpb.GetWalletResponse{Wallet: s.walletMapper(data)}, nil
}

func buildPaginationWallet(meta types.PaginationMeta) *walletpb.Meta {
	return &walletpb.Meta{
		TotalPages:  int32(meta.TotalPages),
		Limit:       int32(meta.Limit),
		TotalDocs:   int32(meta.TotalDocs),
		Page:        int32(meta.Page),
		HasNextPage: meta.HasNextPage,
		// NextPage:    int32(*meta.NextPage),
		// HasPrevPage: meta.HasPrevPage,
	}
}
func (s *server) walletListMapper(data []*model.Wallet) []*walletpb.Wallet {
	var wallets []*walletpb.Wallet
	for _, wallet := range data {
		wallets = append(wallets, s.walletMapper(wallet))
	}
	return wallets
}

// //////////////////////////service Details//////////////
func (s *server) GetAllServices(ctx context.Context, req *servicepb.GetAllServiceDetailsRequest) (*servicepb.GetAllServiceDetailsResponse, error) {
	data, err := s.serviceHandler.GetAllService(ctx, &types.Filter{Page: int(req.Page), PerPage: int(req.PerPage), Search: req.Search})
	if err != nil {
		s.logger.Errorf("Failed to get all wallets: %v", err)
		return nil, err
	}
	return &servicepb.GetAllServiceDetailsResponse{Services: s.serviceListMapper(data.Data), Metadata: buildPaginationService(data.Meta)}, nil
}

func (s *server) GetOneServiceDetail(ctx context.Context, req *servicepb.GetOneServiceDetailRequest) (*servicepb.GetServiceDetailResponse, error) {
	data, err := s.serviceHandler.GetServiceFeeDetail(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("Failed to get wallet: %v", err)
		return nil, err
	}
	return &servicepb.GetServiceDetailResponse{Service: s.MapOneServiceDetail(data)}, nil
}

func buildPaginationService(meta types.PaginationMeta) *servicepb.Meta {
	return &servicepb.Meta{
		TotalPages:  int32(meta.TotalPages),
		Limit:       int32(meta.Limit),
		TotalDocs:   int32(meta.TotalDocs),
		Page:        int32(meta.Page),
		HasNextPage: meta.HasNextPage,
		NextPage:    int32(*meta.NextPage),
		HasPrevPage: meta.HasPrevPage,
	}
}
func (s *server) serviceListMapper(data []*model.ServiceDetails) []*servicepb.ServiceDetails {
	var services []*servicepb.ServiceDetails
	for _, service := range data {
		services = append(services, s.MapServiceDetails(service))
	}
	return services
}
func (s *server) MapServiceDetails(data *model.ServiceDetails) *servicepb.ServiceDetails {
	return &servicepb.ServiceDetails{
		Id:              data.ID.Hex(),
		ServiceCode:     data.ServiceCode,
		ServiceName:     data.ServiceName,
		ServiceType:     data.ServiceType,
		Key:             data.Key,
		AboveAmount:     data.AboveAmount,
		AboveServiceFee: data.AboveServiceFee,
		PaymentType:     data.PaymentType,
		Enabled:         data.Enabled,
		IsDeleted:       data.IsDeleted,
	}
}

// dtoService.ServiceFeeDetailResponse
func (s *server) MapOneServiceDetail(data *dto.ServiceFeeDetailResponse) *servicepb.ServiceDetails {
	return &servicepb.ServiceDetails{
		Id:              data.ID.Hex(),
		ServiceCode:     data.ServiceCode,
		ServiceName:     data.ServiceName,
		ServiceType:     data.ServiceType,
		Key:             data.Key,
		AboveAmount:     data.AboveAmount,
		AboveServiceFee: data.AboveServiceFee,
		PaymentType:     data.PaymentType,
		Enabled:         data.Enabled,
		IsDeleted:       data.IsDeleted,
	}
}

func StartGrpcServer(s *server) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bankpb.RegisterBankServiceServer(grpcServer, s)
	walletpb.RegisterWalletServiceServer(grpcServer, s)
	servicepb.RegisterServiceDetailsServiceServer(grpcServer, s)
	log.Println("gRPC server listening on port 50051")
	return grpcServer, lis
}

func StopGrpcServer(grpcServer *grpc.Server) {
	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}
