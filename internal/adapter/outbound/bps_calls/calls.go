package bpscalls

import (
	"encoding/json"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type BpsCallsInterface interface {
	FetchLinkedAccount(accountNumber string) (*model.LinkedAccount, error)
}
type BpsCalls struct{}

func NewBpsCalls() BpsCallsInterface {
	return &BpsCalls{}
}

func (b *BpsCalls) FetchLinkedAccount(accountNumber string) (*model.LinkedAccount, error) {
	// Simulating an API call to fetch linked account details by account number
	// The End point hasn't been deployed yet, so this is a mock implementation.
	/*
	   url := "https://devcbe.eaglelionsystems.com/api/v1/cbesuperapp/bps/core/account_number"
	   method := "POST"

	   // Create payload
	   payload := strings.NewReader(fmt.Sprintf(`{
	       "account_number": "%s"
	   }`, accountNumber))

	   fmt.Printf("Fetching linked account for account number: %s\n", accountNumber)
	   fmt.Printf("Request URL: %s\n", url)

	   // Create HTTP client and request
	   client := &http.Client{}
	   req, err := http.NewRequest(method, url, payload)
	   if err != nil {
	       fmt.Printf("Error creating request: %v\n", err)
	       return nil, err
	   }

	   // Add headers
	   req.Header.Add("Content-Type", "application/json")
	   req.Header.Add("Accept", "application/json")

	   // Execute request
	   res, err := client.Do(req)
	   if err != nil {
	       fmt.Printf("Error executing request: %v\n", err)
	       return nil, err
	   }
	   defer res.Body.Close()

	   // Check HTTP status code
	   if res.StatusCode != http.StatusOK {
	       fmt.Printf("Unexpected status code: %d\n", res.StatusCode)
	       return nil, err
	   }

	   // Check Content-Type header
	   contentType := res.Header.Get("Content-Type")
	   if !strings.Contains(contentType, "application/json") {
	       fmt.Printf("Unexpected Content-Type: %s\n", contentType)
	       return nil, err
	   }

	   // Read response body
	   body, err := io.ReadAll(res.Body)
	   if err != nil {
	       fmt.Printf("Error reading response body: %v\n", err)
	       return nil, err
	   }

	   // Validate JSON
	   if !json.Valid(body) {
	       fmt.Println("Response is not valid JSON")
	       return nil, err
	   }
	*/

	// Mock data simulating a LinkedAccount response from the API
	// This mock data represents a sample linked account with typical values for testing purposes,
	// including a valid account number, user details, and maker/checker information.
	mockResponse := []byte(`{
        "id": "507f1f77bcf86cd799439011",
        "user_id": "507f191e810c19729de860ea",
        "customer_number": "SIF-f867b034",
        "account_number": "` + accountNumber + `",
        "account_holder_name": "John Doe",
        "account_type": "Savings",
        "branch_code": "BR001",
        "linked_status": true,
        "last_linked_status": false,
        "linked_at": "2025-07-12T11:15:00Z",
        "linker_branch": "Main Branch",
        "registration_type": "Online",
        "is_account_active": true,
        "and_or_status": false,
        "account_branch_code": "BR001",
        "currency": "USD",
        "is_main": true,
        "maker_and_checker": {
            "linkers": {
                "maker": "user123",
                "checker": "admin456"
            },
            "unlinkers": {
                "maker": "",
                "checker": ""
            }
        },
        "created_at": "2025-07-01T10:00:00Z",
        "updated_at": "2025-07-12T11:00:00Z"
    }`)

	// Unmarshal mock response
	var linkedAccount model.LinkedAccount
	if err := json.Unmarshal(mockResponse, &linkedAccount); err != nil {
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	return &linkedAccount, nil
}
