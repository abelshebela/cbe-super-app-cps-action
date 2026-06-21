package customer

import (
	"context"
	"errors"
	"testing"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	accountLookupDTO "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	cps_user_dto "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	account_lookup "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"

	coreio "github.com/hugokessem/coreio/core"
	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ─── approver credentials ──────────────────────────────────────────────────────

const (
	approverID         = "69ea0e540b17af5e50d664f7"
	approverUserCode   = "067879_AB1V_qzcc2m"
	approverFullName   = "DAWIT GIRMA EAGLE"
	approverPhone      = "251936676745"
	approverDepartment = "69ea09868069822215983eaf"
)

func approverCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.ContextKey("user_id"), approverID)
	ctx = context.WithValue(ctx, constants.ContextKey("user_code"), approverUserCode)
	ctx = context.WithValue(ctx, constants.ContextKey("full_name"), approverFullName)
	ctx = context.WithValue(ctx, constants.ContextKey("phone_number"), approverPhone)
	ctx = context.WithValue(ctx, constants.ContextKey("department"), approverDepartment)
	ctx = context.WithValue(ctx, constants.ContextKey("username"), "daveusr")
	return ctx
}

// ─── fixtures ──────────────────────────────────────────────────────────────────

func pendingKYC() *imodel.CustomerKYC {
	return &imodel.CustomerKYC{
		ID:        bson.NewObjectID(),
		UserID:    "user-001",
		KYCStatus: imodel.KYCStatusPending,
		Enabled:   false,
		CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		KYCData: imodel.KYCRequest{
			Sub:               "35725384101447613853835360164180148",
			OriginID:          "357253841014476138538353601641801482",
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
			Address: imodel.Address{
				Woreda: "Woreda 11",
				Zone:   "Addis Ababa",
				Region: "Addis Ababa",
			},
		},
	}
}

func inReviewKYC() *imodel.CustomerKYC {
	k := pendingKYC()
	k.KYCStatus = imodel.KYCStatusInReview
	return k
}

func activeReview(kycID bson.ObjectID) *imodel.StartedKycReview {
	return &imodel.StartedKycReview{
		ID:           bson.NewObjectID(),
		KycID:        kycID,
		ReviewStatus: string(imodel.KYCStatusInReview),
		IsActive:     true,
		StartedAt:    time.Now().Add(-10 * time.Minute),
		ExpiresAt:    time.Now().Add(50 * time.Minute),
	}
}

func approvedCoreResult() *coreio.CusteomerAccountCreationResponse {
	return &coreio.CusteomerAccountCreationResponse{
		CustomerCreationDetail: &coreio.CreateCustomerResult{
			Success: true,
		},
		AccountCreationDetail: &coreio.AccountCreationResult{
			Success:  true,
			Messages: []string{},
		},
	}
}

func cpsActionFor(kyc *imodel.CustomerKYC, requestAction string) *sharedmodel.CPSAction {
	return &sharedmodel.CPSAction{
		UniqueId:      kyc.ID.Hex(),
		RequestAction: requestAction,
		CurrentAction: kyc, // JsonUnmarshal re-marshals the value, so pass the struct directly
	}
}

// ─── mocks ─────────────────────────────────────────────────────────────────────

// mockKYCRepo ──────────────────────────────────────────────

type mockKYCRepo struct {
	kyc             *imodel.CustomerKYC
	review          *imodel.StartedKycReview
	findByIDErr     error
	findReviewErr   error
	createUserErr   error
	updateStatusErr error
	updateReviewErr error
	startReviewErr  error
}

func (m *mockKYCRepo) FindAllWithPagination(_ context.Context, _ types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error) {
	return nil, nil
}
func (m *mockKYCRepo) FindByID(_ context.Context, _ string) (*imodel.CustomerKYC, error) {
	return m.kyc, m.findByIDErr
}
func (m *mockKYCRepo) CreateUser(_ context.Context, _ *coreio.CreateCustomerResult, _ imodel.CustomerKYC) error {
	return m.createUserErr
}
func (m *mockKYCRepo) UpdateKYCStatus(_ context.Context, _, _, _ string, _ bool) error {
	return m.updateStatusErr
}
func (m *mockKYCRepo) FindKycInReview(_ context.Context, _ string) (*imodel.StartedKycReview, error) {
	return m.review, m.findReviewErr
}
func (m *mockKYCRepo) StartKycReview(_ context.Context, r *imodel.StartedKycReview) (*imodel.StartedKycReview, error) {
	return r, m.startReviewErr
}
func (m *mockKYCRepo) UpdateKycReview(_ context.Context, _ string, r *imodel.StartedKycReview) (*imodel.StartedKycReview, error) {
	return r, m.updateReviewErr
}

// mockCPSService ───────────────────────────────────────────

type mockCPSService struct {
	createErr      error
	capturedAction *sharedmodel.CPSAction
}

func (m *mockCPSService) CreateCPSAction(_ context.Context, a *sharedmodel.CPSAction) error {
	m.capturedAction = a
	return m.createErr
}
func (m *mockCPSService) ApproveCPSAction(_ context.Context, _ *sharedmodel.CPSAction) error { return nil }
func (m *mockCPSService) RejectCPSAction(_ context.Context, _ string, _ *sharedmodel.CPSAction) error {
	return nil
}
func (m *mockCPSService) CancelCPSAction(_ context.Context, _ string, _ *sharedmodel.CPSAction) error {
	return nil
}
func (m *mockCPSService) ReverseCPSAction(_ context.Context, _ string) error { return nil }
func (m *mockCPSService) GetCPSActionsByDepartment(_ context.Context, _ string, _ *types.Filter) (*types.PaginatedResponse[[]*sharedmodel.CPSAction], error) {
	return nil, nil
}
func (m *mockCPSService) GetCPSActionsForApprover(_ context.Context, _ string, _ []string, _ *types.Filter) (*types.PaginatedResponse[[]*sharedmodel.CPSAction], string, error) {
	return nil, "", nil
}
func (m *mockCPSService) GetCPSActionsForAuditor(_ context.Context, _ string, _ []string, _ *types.Filter) (*types.PaginatedResponse[[]*sharedmodel.CPSAction], string, error) {
	return nil, "", nil
}
func (m *mockCPSService) GetCPSActions(_ context.Context, _, _ string, _ []string, _ *types.Filter) (*types.PaginatedResponse[[]*sharedmodel.CPSAction], error) {
	return nil, nil
}
func (m *mockCPSService) AuditorClaim(_ context.Context, _ string, _ int) error { return nil }
func (m *mockCPSService) AuditorMark(_ context.Context, _ string, _ sharedmodel.Auditor, _ int) error {
	return nil
}
func (m *mockCPSService) GetUserCreatedActions(_ context.Context, _ string, _ *types.Filter) (*types.PaginatedResponse[[]*sharedmodel.CPSAction], string, error) {
	return nil, "", nil
}
func (m *mockCPSService) GetUserCheckedActions(_ context.Context, _ string, _ *types.Filter) (*types.PaginatedResponse[[]*sharedmodel.CPSAction], string, error) {
	return nil, "", nil
}
func (m *mockCPSService) GetUserAuthorizerIndex(_ context.Context, _ constants.RequestAction) (imodel.CPSActionApproveIndex, error) {
	return imodel.CPSActionApproveIndex{}, nil
}
func (m *mockCPSService) IsMakerOnlyForRequest(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockCPSService) GetActionCountsByDepartemnt(_ context.Context, _ string) (*actionDto.CPSActionCountResponse, error) {
	return nil, nil
}
func (m *mockCPSService) GetCPSActionByID(_ context.Context, _, _ string) (*sharedmodel.CPSAction, error) {
	return nil, nil
}
func (m *mockCPSService) GetCPSActionByUniqueID(_ context.Context, _, _ string) (*sharedmodel.CPSAction, error) {
	return nil, nil
}
func (m *mockCPSService) GetCPSActionByActionCode(_ context.Context, _, _ string) (*sharedmodel.CPSAction, error) {
	return nil, nil
}
func (m *mockCPSService) ExportCpsActionData(_ context.Context, _ []string, _ *types.Filter, _ string) (string, error) {
	return "", nil
}

// mockCPSUserRepo ──────────────────────────────────────────

type mockCPSUserRepo struct {
	user    *imodel.CPSUser
	findErr error
}

func (m *mockCPSUserRepo) FindForExport(_ context.Context, _, _ time.Time, _ string) ([]imodel.ExportCPSUser, error) {
	return nil, nil
}
func (m *mockCPSUserRepo) Create(_ context.Context, _ *imodel.CPSUser) error           { return nil }
func (m *mockCPSUserRepo) Update(_ context.Context, _ string, _ *imodel.CPSUser) error { return nil }
func (m *mockCPSUserRepo) Delete(_ context.Context, _ string) error                     { return nil }
func (m *mockCPSUserRepo) EnableOrDisable(_ context.Context, _ string, _ bool) error   { return nil }
func (m *mockCPSUserRepo) FindByID(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) FindByUserID(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) FindByUsername(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) GetPopulatedByID(_ context.Context, _ string) (*cps_user_dto.CpsUserResponse, error) {
	return nil, nil
}
func (m *mockCPSUserRepo) GetPopulatedWithRole(_ context.Context, _ string) (*cps_user_dto.CpsUserPopulatedResponse, error) {
	return nil, nil
}
func (m *mockCPSUserRepo) GetPopulatedWithRoleByUserName(_ context.Context, _ string) (*cps_user_dto.CpsUserPopulatedResponse, error) {
	return nil, nil
}
func (m *mockCPSUserRepo) FindAllWithPagination(_ context.Context, _ types.Filter) (*types.PaginatedResponse[[]*cps_user_dto.CPSUserWithDepartment], error) {
	return nil, nil
}
func (m *mockCPSUserRepo) FindByPhoneNumber(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) FindByEmail(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) FindByEmailOrPhoneNumberOrUserName(_ context.Context, _, _, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) UpdateCpsUsersJobTitle(_ context.Context, _, _ string) error { return nil }
func (m *mockCPSUserRepo) GetUserByDepartment(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}
func (m *mockCPSUserRepo) GetUserByJobTitle(_ context.Context, _ string) (*imodel.CPSUser, error) {
	return m.user, m.findErr
}

// mockTokenProvider ────────────────────────────────────────

type mockTokenProvider struct {
	token string
	err   error
}

func (m *mockTokenProvider) GetToken(_ context.Context) (string, error) {
	return m.token, m.err
}

// mockCoreAPI ──────────────────────────────────────────────

type mockCoreAPI struct {
	result *coreio.CusteomerAccountCreationResponse
	err    error
}

func (m *mockCoreAPI) AccountCreate(_ context.Context, _ coreio.CreateCustomerParam, _, _ string) (*coreio.CusteomerAccountCreationResponse, error) {
	return m.result, m.err
}
func (m *mockCoreAPI) CreateCustomer(_ context.Context, _ coreio.CreateCustomerParam) (*coreio.CreateCustomerResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) AccountCreation(_ context.Context, _ coreio.AccountCreationParam) (*coreio.AccountCreationResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerLimitFetchByCustomerNumber(_ context.Context, _ coreio.CustomerLimitFetchByCIFParam) (*coreio.CustomerLimitFetchByCIFResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerLimitAmendByCustomerNumber(_ context.Context, _ coreio.CustomerLimitAmendByCIFParam) (*coreio.CustomerLimitAmendByCIFResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerLimitFetchByService(_ context.Context, _ coreio.CustomerLimitFetchByServiceParam) (*coreio.CustomerLimitFetchByServiceResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ServiceLimit(_ context.Context, _ coreio.ServiceLimitParam) (*coreio.ServiceLimitResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) FundTransferVerify(_ context.Context, _ coreio.FundTransferVerifyParam) (*coreio.FundTransferVerifyResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) FundTransfer(_ context.Context, _ coreio.FundTransferParam) (*coreio.FundTransferResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) FundTransferCheck(_ context.Context, _ coreio.FundTransferCheckParam) (*coreio.FundTransferCheckResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) RevertFundTransfer(_ context.Context, _ coreio.RevertFundTransferParam) (*coreio.RevertFundTransferResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) AccountLookup(_ context.Context, _ coreio.AccountLookupParam) (*coreio.AccountLookupResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) LockedAmountFT(_ context.Context, _ coreio.LockedAmountFTParam) (*coreio.LockedAmountFTResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ListLockedAmount(_ context.Context, _ coreio.ListLockedAmountParam) (*coreio.ListLockedAmountResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CreateLockedAmount(_ context.Context, _ coreio.CreateLockedAmountParam) (*coreio.CreateLockedAmountResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ReleaseLockedAmount(_ context.Context, _ coreio.ReleaseLockedAmountParam) (*coreio.ReleaseLockedAmountResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ListStandingOrder(_ context.Context, _ coreio.ListStandingOrderParam) (*coreio.ListStandingOrderResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) UpdateStandingOrder(_ context.Context, _ coreio.UpdateStandingOrderParam) (*coreio.UpdateStandingOrderResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CreateStandingOrder(_ context.Context, _ coreio.CreateStandingOrderParam) (*coreio.CreateStandingOrderResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CancleStandingOrder(_ context.Context, _ coreio.CancleStandingOrderParam) (*coreio.CancelStandingOrderResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ListStandingOrderHistory(_ context.Context, _ coreio.ListStandingOrderHistoryParam) (*coreio.ListStandingOrderHistoryResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) MiniStatementByLimit(_ context.Context, _ coreio.MiniStatementByLimitParams) (*coreio.MiniStatementByLimitResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) MiniStatementByDateRange(_ context.Context, _ coreio.MiniStatementByDateRangeParam) (*coreio.MiniStatementByDateRangeResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) PhoneLookup(_ context.Context, _ coreio.PhoneLookupParam) (*coreio.PhoneLookupResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerLookup(_ context.Context, _ coreio.CustomerLookupParam) (*coreio.CustomerLookupResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) AccountList(_ context.Context, _ coreio.AccountListParam) (*coreio.AccountListResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CardReplace(_ context.Context, _ coreio.CardReplaceParam) (*coreio.CardReplaceResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CardRequest(_ context.Context, _ coreio.CardRequestParam) (*coreio.CardRequestResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) StatusCheck(_ context.Context, _ coreio.StatusCheckParam) (*coreio.StatusCheckResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ExchangeRates(_ context.Context) (*coreio.ExchangeRatesResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerDetail(_ context.Context, _ coreio.CustomerDetailParam) (*coreio.CustomerDetailResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) SplitPayment(_ context.Context, _ coreio.SplitPaymentParam) (*coreio.SplitPaymentResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) ServiceDetail(_ context.Context, _ coreio.ServiceDetailParam) (*coreio.ServiceDetailResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) NameLookup(_ context.Context, _ coreio.NameLookupParam) (*coreio.NameLookupResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) SuperAppSubscribe(_ context.Context, _ coreio.SuperAppSubscribeParam) (*coreio.SuperAppSubscribeResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) SuperAppUnsubscribe(_ context.Context, _ coreio.SuperAppUnsubscribeParam) (*coreio.SuperAppUnsubscribeResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) SuperAppStatusChange(_ context.Context, _ coreio.SuperAppStatusChangeParam) (*coreio.SuperAppStatusChangeResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerLimitFetchByCIFReturnService(_ context.Context, _ coreio.CustomerLimitFetchByCIFReturnServiceParam) (*coreio.CustomerLimitFetchByCIFReturnServiceResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) CustomerFetch(_ context.Context, _ coreio.CustomerFetchParam) (*coreio.CustomerFetchResult, error) {
	return nil, nil
}
func (m *mockCoreAPI) BillPayment(_ context.Context, _ coreio.BillPaymentParam) (*coreio.BillPaymentResult, error) {
	return nil, nil
}

// mockAccountLookup ────────────────────────────────────────

type mockAccountLookup struct{}

func (m *mockAccountLookup) LookupAccountByPhone(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockAccountLookup) LookupAccountByAccountNumber(_ context.Context, _ sharedmodel.AccountLookUpRequest) (*sharedmodel.AccountDetail, error) {
	return nil, nil
}
func (m *mockAccountLookup) LookupAccountByAccountNumberFromBps(_ context.Context, _ string) (accountLookupDTO.AccountResponse, error) {
	return accountLookupDTO.AccountResponse{}, nil
}
func (m *mockAccountLookup) CreateAccountWithFayda(_ context.Context, _ accountLookupDTO.AccountCreateParams) (types.Account, error) {
	return types.Account{}, nil
}
func (m *mockAccountLookup) CifSearch(_ context.Context, _ string) ([]imodel.AccountData, error) {
	return nil, nil
}

var _ account_lookup.Account = (*mockAccountLookup)(nil)

// noopLogger ───────────────────────────────────────────────

type noopLogger struct{}

func (noopLogger) Infof(_ string, _ ...interface{})  {}
func (noopLogger) Warnf(_ string, _ ...interface{})  {}
func (noopLogger) Errorf(_ string, _ ...interface{}) {}
func (noopLogger) Fatalf(_ string, _ ...interface{}) {}
func (noopLogger) Debugf(_ string, _ ...interface{}) {}
func (noopLogger) Sync() error                       { return nil }

// ─── builder ──────────────────────────────────────────────────────────────────

func approverCPSUser() *imodel.CPSUser {
	deptID, _ := bson.ObjectIDFromHex(approverDepartment)
	return &imodel.CPSUser{
		UserCode:    approverUserCode,
		FullName:    approverFullName,
		PhoneNumber: approverPhone,
		Email:       "dawit.eagle@cbe.com.et",
		Department:  deptID,
	}
}

func buildSvc(repo *mockKYCRepo, cps *mockCPSService, userRepo *mockCPSUserRepo, core *mockCoreAPI, token *mockTokenProvider) *customerKYCService {
	return &customerKYCService{
		repo:           repo,
		cpsService:     cps,
		cpsUserRepo:    userRepo,
		accountService: &mockAccountLookup{},
		coreio:         core,
		tokenProvider:  token,
		logger:         noopLogger{},
		cfg: &config.VaultConfig{
			CustomerCreateURL:    "https://superrapp-account-opening-https-ace-uat.apps.cp4itest.cbe.local/cust_creation",
			AccountOpeningURL:    "https://superrapp-account-opening-https-ace-uat.apps.cp4itest.cbe.local/cust_creation",
			CentralKYCBranchCode: "4589",
		},
	}
}

// ─── tests ────────────────────────────────────────────────────────────────────

func TestStartKycReview_HappyPath(t *testing.T) {
	kyc := pendingKYC()
	// FindKycInReview returns "not found" error code → treated as no existing review
	repo := &mockKYCRepo{kyc: kyc, findReviewErr: errors.New("ERROR_RESOURCE_NOT_FOUND")}
	userRepo := &mockCPSUserRepo{user: approverCPSUser()}

	svc := buildSvc(repo, &mockCPSService{}, userRepo, &mockCoreAPI{}, &mockTokenProvider{})
	review, err := svc.StartKycReview(approverCtx(), kyc.ID.Hex())
	if err != nil {
		t.Fatalf("StartKycReview: %v", err)
	}
	// StartKycReview writes directly — no CPS action; it returns the new review record
	if review == nil {
		t.Fatal("expected non-nil review")
	}
	if review.ReviewStatus != string(imodel.KYCStatusInReview) {
		t.Errorf("expected IN_REVIEW status, got %s", review.ReviewStatus)
	}
}

func TestStartKycReview_NotPendingStatus(t *testing.T) {
	kyc := inReviewKYC()
	repo := &mockKYCRepo{kyc: kyc}
	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, &mockTokenProvider{})

	_, err := svc.StartKycReview(approverCtx(), kyc.ID.Hex())
	if err == nil {
		t.Fatal("expected error: only PENDING status allowed")
	}
}

func TestStartKycReview_ActiveReviewAlreadyExists(t *testing.T) {
	kyc := pendingKYC()
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}
	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, &mockTokenProvider{})

	_, err := svc.StartKycReview(approverCtx(), kyc.ID.Hex())
	if err == nil {
		t.Fatal("expected error: active review already exists")
	}
}

func TestEnableOrDisable_Approve_HappyPath(t *testing.T) {
	kyc := inReviewKYC()
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}
	cps := &mockCPSService{}
	svc := buildSvc(repo, cps, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, &mockTokenProvider{})

	err := svc.EnableOrDisable(approverCtx(), kyc.ID.Hex(), "", true)
	if err != nil {
		t.Fatalf("EnableOrDisable(approve): %v", err)
	}
	if cps.capturedAction == nil {
		t.Fatal("expected CPS action to be created")
	}
	if cps.capturedAction.RequestAction != string(constants.RequestApproveCustomerKYC) {
		t.Errorf("wrong action: got %s", cps.capturedAction.RequestAction)
	}
}

func TestEnableOrDisable_Reject_HappyPath(t *testing.T) {
	kyc := inReviewKYC()
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}
	cps := &mockCPSService{}
	svc := buildSvc(repo, cps, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, &mockTokenProvider{})

	err := svc.EnableOrDisable(approverCtx(), kyc.ID.Hex(), "document mismatch", false)
	if err != nil {
		t.Fatalf("EnableOrDisable(reject): %v", err)
	}
	if cps.capturedAction.RequestAction != string(constants.RequestRejectCustomerKYC) {
		t.Errorf("wrong action: got %s", cps.capturedAction.RequestAction)
	}
}

func TestEnableOrDisable_MissingUserContext(t *testing.T) {
	kyc := inReviewKYC()
	repo := &mockKYCRepo{kyc: kyc}
	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{}, &mockCoreAPI{}, &mockTokenProvider{})

	err := svc.EnableOrDisable(context.Background(), kyc.ID.Hex(), "", true)
	if err == nil {
		t.Fatal("expected error: missing user context")
	}
}

func TestAuthorize_Approve_HappyPath(t *testing.T) {
	kyc := inReviewKYC()
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}
	token := &mockTokenProvider{token: "test-bearer-token"}
	core := &mockCoreAPI{result: approvedCoreResult()}

	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, core, token)
	action := cpsActionFor(kyc, string(constants.RequestApproveCustomerKYC))

	result, err := svc.Authorize(approverCtx(), action)
	if err != nil {
		t.Fatalf("Authorize(approve): %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil CPSAction result")
	}
}

func TestAuthorize_Approve_TokenError(t *testing.T) {
	kyc := inReviewKYC()
	repo := &mockKYCRepo{kyc: kyc}
	token := &mockTokenProvider{err: errors.New("token service unavailable")}

	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, token)
	action := cpsActionFor(kyc, string(constants.RequestApproveCustomerKYC))

	_, err := svc.Authorize(approverCtx(), action)
	if err == nil {
		t.Fatal("expected error when token fetch fails")
	}
}

func TestAuthorize_Approve_T24Error(t *testing.T) {
	kyc := inReviewKYC()
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}
	token := &mockTokenProvider{token: "test-token"}
	core := &mockCoreAPI{err: errors.New("T24 SOAP: LEGAL.ID TOO MANY CHARACTERS")}

	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, core, token)
	action := cpsActionFor(kyc, string(constants.RequestApproveCustomerKYC))

	_, err := svc.Authorize(approverCtx(), action)
	if err == nil {
		t.Fatal("expected error when T24 call fails")
	}
}

func TestAuthorize_Reject_HappyPath(t *testing.T) {
	kyc := inReviewKYC()
	kyc.KYCRejectReason = "identity document expired"
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}

	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, &mockTokenProvider{})
	action := cpsActionFor(kyc, string(constants.RequestRejectCustomerKYC))

	result, err := svc.Authorize(approverCtx(), action)
	if err != nil {
		t.Fatalf("Authorize(reject): %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestAuthorize_Reject_MissingReason(t *testing.T) {
	kyc := inReviewKYC()
	// KYCRejectReason intentionally empty
	review := activeReview(kyc.ID)
	repo := &mockKYCRepo{kyc: kyc, review: review}

	svc := buildSvc(repo, &mockCPSService{}, &mockCPSUserRepo{user: approverCPSUser()}, &mockCoreAPI{}, &mockTokenProvider{})
	action := cpsActionFor(kyc, string(constants.RequestRejectCustomerKYC))

	_, err := svc.Authorize(approverCtx(), action)
	if err == nil {
		t.Fatal("expected error: rejection reason is required")
	}
}
