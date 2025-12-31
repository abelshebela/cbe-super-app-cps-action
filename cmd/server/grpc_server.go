package server

import (
	bankpb "cbe-super-app-cps-action/grpc/bank"
	servicepb "cbe-super-app-cps-action/grpc/service/proto"

	topuppb "cbe-super-app-cps-action/grpc/topup/proto"
	walletpb "cbe-super-app-cps-action/grpc/wallet/proto"
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"net"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server implements bankpb.BankServiceServer
// Embed UnimplementedBankServiceServer for forward compatibility

type server struct {
	bankpb.UnimplementedBankServiceServer
	walletpb.UnimplementedWalletServiceServer
	servicepb.UnimplementedServiceDetailsServiceServer
	topuppb.UnimplementedTopupServiceServer
	bankHandler    service.BankService
	walletHandler  service.WalletService
	serviceHandler service.ServicesService
	topupHandler   service.TopupService
	logger         utils.Logger
}

func NewGrpcServer(bankHandler service.BankService, walletHandler service.WalletService, serviceHandler service.ServicesService, topupHandler service.TopupService, logger utils.Logger) *server {
	return &server{
		bankHandler:    bankHandler,
		walletHandler:  walletHandler,
		serviceHandler: serviceHandler,
		topupHandler:   topupHandler,
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

func (s *server) walletMapper(data *local_model.Wallet) *walletpb.Wallet {
	return &walletpb.Wallet{
		Id:         data.ID.Hex(),
		Name:       data.Name,
		Avatar:     data.Avatar,
		UniqueCode: data.UniqueCode,
		IsDeleted:  data.IsDeleted,
		Enabled:    data.Enabled,
		Services: &walletpb.Services{
			Self:  data.Services.Self,
			Other: data.Services.Other,
			Agent: data.Services.Agent,
		},
	}
}

func (s *server) bankListMapper(data []model.Bank) []*bankpb.Bank {
	var banks []*bankpb.Bank
	for i := range data {
		banks = append(banks, s.bankMapper(&data[i]))
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
		BicCode: data.BICCode,
		Logo:    data.Logo,
		Enabled: data.Enabled,
		Type:    data.Type,
	}
}
func buildPagination(meta types.PaginationMeta) *bankpb.Meta {
	var nextPage int32
	if meta.NextPage != nil {
		nextPage = int32(*meta.NextPage)
	}
	return &bankpb.Meta{
		TotalPages:  int32(meta.TotalPages),
		Limit:       int32(meta.Limit),
		TotalDocs:   int32(meta.TotalDocs),
		Page:        int32(meta.Page),
		HasNextPage: meta.HasNextPage,
		NextPage:    nextPage,
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
func (s *server) walletListMapper(data []local_model.Wallet) []*walletpb.Wallet {
	var wallets []*walletpb.Wallet
	for i := range data {
		wallets = append(wallets, s.walletMapper(&data[i]))
	}
	return wallets
}

// //////////////////////////service Details//////////////
// func (s *server) GetAllServices(ctx context.Context, req *servicepb.GetAllServiceDetailsRequest) (*servicepb.GetAllServiceDetailsResponse, error) {
// 	data, err := s.serviceHandler.GetAllService(ctx, &types.Filter{Page: int(req.Page), PerPage: int(req.PerPage), Search: req.Search})
// 	if err != nil {
// 		s.logger.Errorf("Failed to get all wallets: %v", err)
// 		return nil, err
// 	}
// 	return &servicepb.GetAllServiceDetailsResponse{Services: s.serviceListMapper(data.Data), Metadata: buildPaginationService(data.Meta)}, nil
// }

// func (s *server) GetOneServiceDetail(ctx context.Context, req *servicepb.GetOneServiceDetailRequest) (*servicepb.GetServiceDetailResponse, error) {
// 	data, err := s.serviceHandler.GetServiceFeeDetail(ctx, req.Id)
// 	if err != nil {
// 		s.logger.Errorf("Failed to get wallet: %v", err)
// 		return nil, err
// 	}
// 	return &servicepb.GetServiceDetailResponse{Service: s.MapOneServiceDetail(data)}, nil
// }

func buildPaginationService(meta types.PaginationMeta) *servicepb.Meta {
	var nextPage int32
	if meta.NextPage != nil {
		nextPage = int32(*meta.NextPage)
	}
	return &servicepb.Meta{
		TotalPages:  int32(meta.TotalPages),
		Limit:       int32(meta.Limit),
		TotalDocs:   int32(meta.TotalDocs),
		Page:        int32(meta.Page),
		HasNextPage: meta.HasNextPage,
		NextPage:    nextPage,
		HasPrevPage: meta.HasPrevPage,
	}
}
func (s *server) serviceListMapper(data []imodel.ServiceDetails) []*servicepb.ServiceDetails {
	var services []*servicepb.ServiceDetails
	for i := range data {
		services = append(services, s.MapServiceDetails(&data[i]))
	}
	return services
}
func (s *server) MapServiceDetails(data *imodel.ServiceDetails) *servicepb.ServiceDetails {
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

func (s *server) GetAllTopup(ctx context.Context, req *topuppb.TopupRequest) (*topuppb.TopupResponse, error) {
	data, err := s.topupHandler.GetAllTopup(ctx, types.Filter{})
	if err != nil {
		s.logger.Errorf("Failed to get topup: %v", err)
		return nil, err
	}
	return &topuppb.TopupResponse{Topups: s.TopupMapper(data.Data)}, nil
}

func (s *server) TopupMapper(data []model.Topup) []*topuppb.Topup {
	var topups []*topuppb.Topup
	for i := range data {
		topups = append(topups, &topuppb.Topup{
			Id:             data[i].ID.Hex(),
			Name:           data[i].Name,
			Code:           data[i].Code,
			Avatar:         data[i].Avatar,
			Enabled:        data[i].Enabled,
			IsDeleted:      data[i].IsDeleted,
			CreatedAt:      timestamppb.New(data[i].CreatedAt),
			LastModifiedAt: timestamppb.New(data[i].LastModifiedAt),
			DeletedAt:      timestamppb.New(data[i].DeletedAt),
		})
	}
	return topups
}

func StartGrpcServer(s *server, logger utils.Logger) (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	// create gRPC server with OpenTelemetry stats handler
	statsHandler := otelgrpc.NewServerHandler()
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(statsHandler),
	)
	bankpb.RegisterBankServiceServer(grpcServer, s)
	walletpb.RegisterWalletServiceServer(grpcServer, s)
	servicepb.RegisterServiceDetailsServiceServer(grpcServer, s)
	topuppb.RegisterTopupServiceServer(grpcServer, s)
	logger.Infof("gRPC server listening on port 50051")
	return grpcServer, lis
}

func StopGrpcServer(grpcServer *grpc.Server, logger utils.Logger) {
	logger.Infof("Stopping gRPC server...")
	grpcServer.GracefulStop()
	logger.Infof("gRPC server stopped")
}
