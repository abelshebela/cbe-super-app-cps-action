package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Bank represents the bank data structure
type Bank struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Logo           string    `json:"logo"`
	Code           string    `json:"code"`
	Bic            string    `json:"bic"`
	Enabled        bool      `json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

// Request/Response types (matching the server types)
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

// BankServiceClient interface for the gRPC client
type BankServiceClient interface {
	GetBank(ctx context.Context, req *GetBankRequest) (*GetBankResponse, error)
	GetAllBanks(ctx context.Context, req *GetAllBanksRequest) (*GetAllBanksResponse, error)
	CreateBank(ctx context.Context, req *CreateBankRequest) (*CreateBankResponse, error)
}

// bankServiceClient implements BankServiceClient
type bankServiceClient struct {
	conn *grpc.ClientConn
}

// NewBankServiceClient creates a new bank service client
func NewBankServiceClient(conn *grpc.ClientConn) BankServiceClient {
	return &bankServiceClient{conn: conn}
}

// For this example, we'll simulate the gRPC calls using JSON over HTTP
// In a real implementation, these would be actual gRPC calls
func (c *bankServiceClient) GetBank(ctx context.Context, req *GetBankRequest) (*GetBankResponse, error) {
	// This is a placeholder implementation
	// In reality, this would make a gRPC call to the server
	fmt.Printf("Making gRPC call: GetBank with ID: %s\n", req.Id)
	
	// Simulate response
	return &GetBankResponse{
		Bank: &Bank{
			Id:      req.Id,
			Name:    "Example Bank",
			Code:    "EXB",
			Bic:     "EXBKET22",
			Enabled: true,
		},
		Message: "Bank retrieved successfully",
	}, nil
}

func (c *bankServiceClient) GetAllBanks(ctx context.Context, req *GetAllBanksRequest) (*GetAllBanksResponse, error) {
	fmt.Printf("Making gRPC call: GetAllBanks with page: %d, per_page: %d\n", req.Page, req.PerPage)
	
	// Simulate response
	return &GetAllBanksResponse{
		Banks: []*Bank{
			{
				Id:      "1",
				Name:    "Commercial Bank of Ethiopia",
				Code:    "CBE",
				Bic:     "CBEKET22",
				Enabled: true,
			},
			{
				Id:      "2",
				Name:    "Awash Bank",
				Code:    "AWB",
				Bic:     "AWBKET22",
				Enabled: true,
			},
		},
		Total:      2,
		Page:       req.Page,
		PerPage:    req.PerPage,
		TotalPages: 1,
		Message:    "Banks retrieved successfully",
	}, nil
}

func (c *bankServiceClient) CreateBank(ctx context.Context, req *CreateBankRequest) (*CreateBankResponse, error) {
	fmt.Printf("Making gRPC call: CreateBank with name: %s, code: %s\n", req.Name, req.Code)
	
	// Simulate response
	return &CreateBankResponse{
		Message:    "Bank creation request submitted successfully",
		ActionCode: "ACT-" + fmt.Sprintf("%d", time.Now().Unix()),
	}, nil
}

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create client
	client := NewBankServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== gRPC Bank Service Client Example ===\n")

	// Test GetAllBanks
	fmt.Println("1. Testing GetAllBanks:")
	getAllReq := &GetAllBanksRequest{
		Page:    1,
		PerPage: 10,
	}
	
	getAllResp, err := client.GetAllBanks(ctx, getAllReq)
	if err != nil {
		log.Printf("GetAllBanks failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", formatJSON(getAllResp))
	}
	fmt.Println()

	// Test GetBank
	fmt.Println("2. Testing GetBank:")
	getBankReq := &GetBankRequest{
		Id: "1",
	}
	
	getBankResp, err := client.GetBank(ctx, getBankReq)
	if err != nil {
		log.Printf("GetBank failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", formatJSON(getBankResp))
	}
	fmt.Println()

	// Test CreateBank
	fmt.Println("3. Testing CreateBank:")
	createReq := &CreateBankRequest{
		Name: "Test Bank",
		Code: "TSB",
		Bic:  "TSBKET22",
		Logo: []byte("dummy_logo_data"),
	}
	
	createResp, err := client.CreateBank(ctx, createReq)
	if err != nil {
		log.Printf("CreateBank failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", formatJSON(createResp))
	}
	fmt.Println()

	fmt.Println("=== gRPC Client Example Completed ===")
}

// formatJSON formats a struct as pretty JSON
func formatJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(b)
}
