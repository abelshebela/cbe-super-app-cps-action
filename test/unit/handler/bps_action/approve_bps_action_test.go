package bps_action_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cbe-super-app-cps-action/internal/constants"
	bpsActionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/bps_action_handler"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/test/testinit"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeBPSActionService struct {
	service.BPSActionService
	GetByCodeFn func(ctx context.Context, uniqueID, department string) (*bps.BPSAction, error)
	ApproveFn   func(ctx context.Context, action *bps.BPSAction) error
}

func (f *fakeBPSActionService) ApproveBPSAction(ctx context.Context, action *bps.BPSAction) error {
	if f.ApproveFn != nil {
		return f.ApproveFn(ctx, action)
	}
	return nil
}

func (f *fakeBPSActionService) GetBPSActionByActionCode(ctx context.Context, uniqueID, department string) (*bps.BPSAction, error) {
	if f.GetByCodeFn != nil {
		return f.GetByCodeFn(ctx, uniqueID, department)
	}
	return &bps.BPSAction{}, nil
}

func (f *fakeBPSActionService) RejectBPSAction(ctx context.Context, actionCode string, action *bps.BPSAction) error {
	_ = actionCode
	return nil
}

func (f *fakeBPSActionService) GetBPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps.BPSAction], error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetBPSActionsForApprover(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps.BPSAction], error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetBPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps.BPSAction], error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetBPSActions(ctx context.Context, userID, role string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps.BPSAction], error) {
	return nil, nil
}

func (f *fakeBPSActionService) AuditorClaim(ctx context.Context, actionCode string, activeGroup int) error {
	return nil
}

func (f *fakeBPSActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	return nil
}

func (f *fakeBPSActionService) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps.BPSAction], error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps.BPSAction], error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetUserAuthorizerIndex(ctx context.Context, requestAction constants.RequestAction) (imodel.BPSActionApproveIndex, error) {
	return imodel.BPSActionApproveIndex{}, nil
}

func (f *fakeBPSActionService) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	return false, nil
}

func (f *fakeBPSActionService) GetActionCountsByDepartemnt(ctx context.Context, department string) (*bpsActionDto.BPSActionCountResponse, error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetBPSActionByID(ctx context.Context, id, department string) (*bps.BPSAction, error) {
	return nil, nil
}

func (f *fakeBPSActionService) GetBPSActionByUniqueID(ctx context.Context, id, department string) (*bps.BPSAction, error) {
	return nil, nil
}

func TestApproveBPSActionOK(t *testing.T) {
	logger := testinit.NewLogger()

	fakeSvc := &fakeBPSActionService{
		GetByCodeFn: func(ctx context.Context, uniqueID, department string) (*bps.BPSAction, error) {
			return &bps.BPSAction{ActionCode: uniqueID, ID: bson.NewObjectID()}, nil
		},
		ApproveFn: func(ctx context.Context, action *bps.BPSAction) error { return nil },
	}

	adapter := bps_action_handler.InitBPSActionAdapter(fakeSvc, logger)

	req := httptest.NewRequest(http.MethodPatch, "/actions/AC123/approve", nil)
	req = testinit.WithUserContext(req, "user-1", "Test User", "251900000000", "IT", "CHECKER", "U001")
	req = testinit.WithRoleCode(req, "ROLE_1")
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("action_code", "AC123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, "req-1"))

	rr := httptest.NewRecorder()

	adapter.ApproveBPSAction(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp localization.StandardResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v. body=%s", err, rr.Body.String())
	}
	if !resp.Ok {
		t.Fatalf("expected ok=true, got ok=false. body=%s", rr.Body.String())
	}
}
