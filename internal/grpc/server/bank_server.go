package server

import (
	"context"
	"mime/multipart"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// BankServer implements the gRPC bank service
type BankServer struct {
	bankHandler bank.BankHandlerService
	logger      utils.Logger
}

// NewBankServer creates a new bank gRPC server
func NewBankServer(bankHandler bank.BankHandlerService, logger utils.Logger) *BankServer {
	return &BankServer{
		bankHandler: bankHandler,
		logger:      logger,
	}
}

// Bank protobuf message
type Bank struct {
	Id             string                 `json:"id"`
	Name           string                 `json:"name"`
	Logo           string                 `json:"logo"`
	Code           string                 `json:"code"`
	Bic            string                 `json:"bic"`
	Enabled        bool                   `json:"enabled"`
	CreatedAt      *timestamppb.Timestamp `json:"created_at"`
	LastModifiedAt *timestamppb.Timestamp `json:"last_modified_at"`
}

// Request/Response types
type GetBankRequest struct {
	Id string `json:"id"`
}

type GetBankResponse struct {
	Bank    *Bank  `json:"bank"`
	Message string `json:"message"`
}

type GetAllBanksRequest struct {
	Page      int32  `json:"page"`
	PerPage   int32  `json:"per_page"`
	Search    string `json:"search"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type GetAllBanksResponse struct {
	Banks      []*Bank `json:"banks"`
	Total      int32   `json:"total"`
	Page       int32   `json:"page"`
	PerPage    int32   `json:"per_page"`
	TotalPages int32   `json:"total_pages"`
	Message    string  `json:"message"`
}

type CreateBankRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Bic  string `json:"bic"`
	Logo []byte `json:"logo"`
}

type CreateBankResponse struct {
	Message    string `json:"message"`
	ActionCode string `json:"action_code"`
}

type UpdateBankRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Bic  string `json:"bic"`
}

type UpdateBankResponse struct {
	Message    string `json:"message"`
	ActionCode string `json:"action_code"`
}

type DeleteBankRequest struct {
	Id string `json:"id"`
}

type DeleteBankResponse struct {
	Message    string `json:"message"`
	ActionCode string `json:"action_code"`
}

type EnableDisableBankRequest struct {
	Id     string `json:"id"`
	Enable bool   `json:"enable"`
}

type EnableDisableBankResponse struct {
	Message    string `json:"message"`
	ActionCode string `json:"action_code"`
}

type UpdateBankLogoRequest struct {
	Id   string `json:"id"`
	Logo []byte `json:"logo"`
}

type UpdateBankLogoResponse struct {
	Message    string `json:"message"`
	ActionCode string `json:"action_code"`
}

// Helper function to convert entity.Bank to protobuf Bank
func (s *BankServer) convertEntityToProto(bank *entity.Bank) *Bank {
	return &Bank{
		Id:             bank.ID,
		Name:           bank.Name,
		Logo:           bank.Logo,
		Code:           bank.Code,
		Bic:            bank.BIC,
		Enabled:        bank.Enabled,
		CreatedAt:      timestamppb.New(bank.CreatedAt),
		LastModifiedAt: timestamppb.New(bank.LastModifiedAt),
	}
}

// Helper function to create CPS user context (simplified for gRPC)
func (s *BankServer) createCPSUserForCreate() *model.CreateCPSAction {
	return &model.CreateCPSAction{
		MakerUser: model.User{
			UserCode:    "grpc-user",
			FullName:    "gRPC Client",
			PhoneNumber: "000000000",
			Department:  "API",
		},
		Department: "API",
	}
}

// GetBank retrieves a single bank by ID
func (s *BankServer) GetBank(ctx context.Context, req *GetBankRequest) (*GetBankResponse, error) {
	s.logger.Infof("gRPC GetBank called with ID: %s", req.Id)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "bank ID is required")
	}

	bank, err := s.bankHandler.GetOneBank(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("Failed to get bank: %v", err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetBankResponse{
		Bank:    s.convertEntityToProto(bank),
		Message: "Bank retrieved successfully",
	}, nil
}

// GetAllBanks retrieves all banks with pagination
func (s *BankServer) GetAllBanks(ctx context.Context, req *GetAllBanksRequest) (*GetAllBanksResponse, error) {
	s.logger.Infof("gRPC GetAllBanks called")

	// Set default pagination if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	filterParams := &constant.MongoFilter{
		Page:    int(req.Page),
		PerPage: int(req.PerPage),
		Search:  req.Search,
		Filters: make(map[string]interface{}),
	}

	result, err := s.bankHandler.GetAllBank(ctx, filterParams)
	if err != nil {
		s.logger.Errorf("Failed to get all banks: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert entities to protobuf messages
	protoBanks := make([]*Bank, len(result.Data))
	for i, bank := range result.Data {
		protoBanks[i] = s.convertEntityToProto(bank)
	}

	return &GetAllBanksResponse{
		Banks:      protoBanks,
		Total:      int32(result.Meta.TotalDocs),
		Page:       int32(result.Meta.Page),
		PerPage:    int32(result.Meta.Limit),
		TotalPages: int32(result.Meta.TotalPages),
		Message:    "Banks retrieved successfully",
	}, nil
}

// CreateBank creates a new bank
func (s *BankServer) CreateBank(ctx context.Context, req *CreateBankRequest) (*CreateBankResponse, error) {
	s.logger.Infof("gRPC CreateBank called for bank: %s", req.Name)

	if req.Name == "" || req.Code == "" || req.Bic == "" {
		return nil, status.Error(codes.InvalidArgument, "name, code, and BIC are required")
	}

	// Create a mock file header for logo (simplified for gRPC)
	var logoHeader *multipart.FileHeader
	if len(req.Logo) > 0 {
		logoHeader = &multipart.FileHeader{
			Filename: "logo.png",
			Size:     int64(len(req.Logo)),
		}
	}

	bankRequest := dto.CreateBankRequest{
		Name: req.Name,
		Code: req.Code,
		BIC:  req.Bic,
		Logo: logoHeader,
	}

	cpsRequest := s.createCPSUserForCreate()
	cpsRequest.ActionData = bankRequest

	result, err := s.bankHandler.CreateOneBank(ctx, *cpsRequest)
	if err != nil {
		s.logger.Errorf("Failed to create bank: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateBankResponse{
		Message:    "Bank creation request submitted successfully",
		ActionCode: result.ActionCode,
	}, nil
}

// UpdateBank updates an existing bank
func (s *BankServer) UpdateBank(ctx context.Context, req *UpdateBankRequest) (*UpdateBankResponse, error) {
	s.logger.Infof("gRPC UpdateBank called for bank ID: %s", req.Id)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "bank ID is required")
	}

	updateRequest := dto.UpdateBankRequest{
		ID:   req.Id,
		Name: req.Name,
		Code: req.Code,
		BIC:  req.Bic,
	}

	cpsRequest := s.createCPSUserForCreate()
	cpsRequest.ActionData = updateRequest

	result, err := s.bankHandler.UpdateOneBank(ctx, req.Id, *cpsRequest)
	if err != nil {
		s.logger.Errorf("Failed to update bank: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &UpdateBankResponse{
		Message:    "Bank update request submitted successfully",
		ActionCode: result.ActionCode,
	}, nil
}

// DeleteBank deletes a bank
func (s *BankServer) DeleteBank(ctx context.Context, req *DeleteBankRequest) (*DeleteBankResponse, error) {
	s.logger.Infof("gRPC DeleteBank called for bank ID: %s", req.Id)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "bank ID is required")
	}

	cpsRequest := s.createCPSUserForCreate()
	cpsRequest.ActionData = entity.Bank{ID: req.Id}

	result, err := s.bankHandler.DeleteOneBank(ctx, req.Id, *cpsRequest)
	if err != nil {
		s.logger.Errorf("Failed to delete bank: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &DeleteBankResponse{
		Message:    "Bank deletion request submitted successfully",
		ActionCode: result.ActionCode,
	}, nil
}

// EnableDisableBank enables or disables a bank
func (s *BankServer) EnableDisableBank(ctx context.Context, req *EnableDisableBankRequest) (*EnableDisableBankResponse, error) {
	s.logger.Infof("gRPC EnableDisableBank called for bank ID: %s, enable: %v", req.Id, req.Enable)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "bank ID is required")
	}

	var requestAction model.RequestAction
	if req.Enable {
		requestAction = model.RequestEnableBank
	} else {
		requestAction = model.RequestDisableBank
	}

	cpsRequest := s.createCPSUserForCreate()
	cpsRequest.ActionData = dto.UpdateBankRequest{ID: req.Id}

	result, err := s.bankHandler.EnableOrDisableBank(ctx, req.Id, requestAction, *cpsRequest)
	if err != nil {
		s.logger.Errorf("Failed to enable/disable bank: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	action := "disable"
	if req.Enable {
		action = "enable"
	}

	return &EnableDisableBankResponse{
		Message:    "Bank " + action + " request submitted successfully",
		ActionCode: result.ActionCode,
	}, nil
}

// UpdateBankLogo updates a bank's logo
func (s *BankServer) UpdateBankLogo(ctx context.Context, req *UpdateBankLogoRequest) (*UpdateBankLogoResponse, error) {
	s.logger.Infof("gRPC UpdateBankLogo called for bank ID: %s", req.Id)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "bank ID is required")
	}

	if len(req.Logo) == 0 {
		return nil, status.Error(codes.InvalidArgument, "logo data is required")
	}

	// Create a mock file header for logo
	logoHeader := &multipart.FileHeader{
		Filename: "logo.png",
		Size:     int64(len(req.Logo)),
	}

	updateLogo := dto.UpdateLogo{
		ID:   req.Id,
		Logo: logoHeader,
	}

	cpsRequest := s.createCPSUserForCreate()
	cpsRequest.ActionData = updateLogo

	result, err := s.bankHandler.UpdateLogo(ctx, req.Id, *cpsRequest)
	if err != nil {
		s.logger.Errorf("Failed to update bank logo: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &UpdateBankLogoResponse{
		Message:    "Bank logo update request submitted successfully",
		ActionCode: result.ActionCode,
	}, nil
}
