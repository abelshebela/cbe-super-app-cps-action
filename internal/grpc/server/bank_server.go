package server

import (
	"context"
	"log"
	"net"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	bankpb "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/grpc/bank"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"google.golang.org/grpc"
)

// server implements bankpb.BankServiceServer
// Embed UnimplementedBankServiceServer for forward compatibility

type server struct {
	bankpb.UnimplementedBankServiceServer
	bankHandler bank.BankHandlerService
	logger      utils.Logger
}

func NewGrpcBankServer(handler bank.BankHandlerService, logger utils.Logger) bankpb.BankServiceServer {
	server := server{
		bankHandler: handler,
		logger:      logger,
	}

	return &server
}
func (s *server) GetOneBank(ctx context.Context, req *bankpb.GetOneBankRequest) (*bankpb.GetOneBankResponse, error) {
	// TODO: Implement logic
	data, err := s.bankHandler.GetOneBank(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("Failed to get bank: %v", err)
		return nil, err
	}
	return &bankpb.GetOneBankResponse{Bank: s.bankMapper(data)}, nil
}

func (s *server) bankMapper(data *entity.Bank) *bankpb.Bank {
	return &bankpb.Bank{
		Id:      data.ID,
		Name:    data.Name,
		Bic:     data.BIC,
		Code:    data.Code,
		Logo:    data.Logo,
		Enabled: data.Enabled,
	}
}

func (s *server) bankListMapper(data []*entity.Bank) []*bankpb.Bank {
	var banks []*bankpb.Bank
	for _, bank := range data {
		banks = append(banks, s.bankMapper(bank))
	}
	return banks
}

func (s *server) GetAllBank(ctx context.Context, req *bankpb.GetAllBankRequest) (*bankpb.GetAllBankResponse, error) {

	data, err := s.bankHandler.GetAllBank(ctx, &constant.MongoFilter{Page: int(req.Page), PerPage: int(req.PerPage), Search: req.Search})
	if err != nil {
		s.logger.Errorf("Failed to get all banks: %v", err)
		return nil, err
	}

	return &bankpb.GetAllBankResponse{Banks: s.bankListMapper(data.Data), Metadata: buildPagination(data.Meta)}, nil
}

func buildPagination(meta common_util.PaginationMeta) *bankpb.Meta {
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

func StartGrpcServer(server bankpb.BankServiceServer) {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bankpb.RegisterBankServiceServer(grpcServer, server)
	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StopGrpcServer(server bankpb.BankServiceServer) {
	grpcServer := grpc.NewServer()
	bankpb.RegisterBankServiceServer(grpcServer, server)
	log.Println("gRPC server stopped")
	grpcServer.Stop()
}
