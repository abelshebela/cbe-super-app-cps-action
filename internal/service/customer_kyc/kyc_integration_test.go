package customer

import (
	"context"
	"os"
	"testing"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	account_lookup "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	token_provider_svc "cbe-super-app-cps-action/internal/service/token_provider"
	tp_client "cbe-super-app-cps-action/internal/storage/external_call/token_provider"

	coreio "github.com/hugokessem/coreio/core"
	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Run with:
//
//	go test ./internal/service/customer_kyc/... -v -run TestIntegration_Authorize_RealAccountCreation
//
// Requires VAULT_ADDR, VAULT_TOKEN, VAULT_PATH env vars (see .env).
func TestIntegration_Authorize_RealAccountCreation(t *testing.T) {
	if os.Getenv("VAULT_ADDR") == "" {
		t.Skip("VAULT_ADDR not set — load .env before running integration tests")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("config.Load failed (Vault unreachable or misconfigured): %v", err)
	}

	logger := noopLogger{}

	tokenClient := tp_client.NewTokenProviderClient(
		cfg.AccountOpeningTokenURL,
		cfg.AccountOpeningTokenClientID,
		cfg.AccountOpeningTokenClientSecret,
		cfg.AccountOpeningTokenScope,
		logger,
	)
	tokenSvc := token_provider_svc.NewTokenProviderService(tokenClient, nil, logger)

	coreAPI := coreio.NewCBECoreAPI(&coreio.CBECoreCredential{
		Username: cfg.CbeCoreUsername,
		Password: cfg.CbeCorePassword,
		Url:      cfg.CbeCoreUrl,
	})

	kyc := &imodel.CustomerKYC{
		ID:        bson.NewObjectID(),
		KYCStatus: imodel.KYCStatusInReview,
		KYCData: imodel.KYCRequest{
			FullName:          "MEREWA SELH KEDIR",
			PhoneNumber:       "0986967210",
			Gender:            "FEMALE",
			BirthDate:         time.Date(1985, 12, 1, 0, 0, 0, 0, time.UTC),
			Nationality:       "ET",
			Country:           "ET",
			EmployementStatus: "EMPLOYED",
			MonthlyIncome:     "123450",
			Currency:          "ETB",
			SubAccountType:    "RETAIL",
			MothersName:       "MAMI LEMA LETA",
			Vendor:            imodel.Fayda,
			IssuedDate:        "20210101",
			OriginID:          "35725384101447613853835360",
			Sub:               "35725384101447613853835360",
			Address: imodel.Address{
				Woreda: "Woreda 11",
				Zone:   "Addis Ababa",
				Region: "Addis Ababa",
			},
		},
	}

	svc := &customerKYCService{
		repo:           &mockKYCRepo{kyc: kyc, review: activeReview(kyc.ID)},
		cpsService:     &mockCPSService{},
		cpsUserRepo:    &mockCPSUserRepo{user: approverCPSUser()},
		accountService: account_lookup.NewCoreAccountLookupAdapter(coreAPI, "", 0, logger),
		coreio:         coreAPI,
		tokenProvider:  tokenSvc,
		logger:         logger,
		cfg:            cfg,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.ContextKey("user_id"), approverID)
	ctx = context.WithValue(ctx, constants.ContextKey("user_code"), approverUserCode)
	ctx = context.WithValue(ctx, constants.ContextKey("full_name"), approverFullName)
	ctx = context.WithValue(ctx, constants.ContextKey("phone_number"), approverPhone)
	ctx = context.WithValue(ctx, constants.ContextKey("department"), approverDepartment)

	action := &sharedmodel.CPSAction{
		UniqueId:      kyc.ID.Hex(),
		RequestAction: string(constants.RequestApproveCustomerKYC),
		CurrentAction: kyc,
	}

	result, err := svc.Authorize(ctx, action)
	if err != nil {
		t.Fatalf("Authorize (real T24): %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil CPSAction result")
	}
	t.Logf("T24 account creation succeeded: %+v", result)
}
