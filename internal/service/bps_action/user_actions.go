package bps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Ensure this file compiles within the cpsaction package and uses the existing service struct
var _ = utils.Logger(nil)
var _ = storage.CPSActionRepository(nil)

// GetUserApprovedCPSActions returns approved CPS actions created by a specific user (maker)
func (ca *bpsActionService) GetUserApprovedCPSActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserApprovedCPSActions", "CPSAction", "GetUserApprovedCPSActions")
	defer span.End()

	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID
	filterParams.Filters["action_status"] = string(constants.Approved)

	result, err := ba.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find approved cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

// GetUserPendingCPSActions returns pending CPS actions created by a specific user (maker)
func (ca *bpsActionService) GetUserPendingCPSActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserPendingCPSActions", "CPSAction", "GetUserPendingCPSActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID
	filterParams.Filters["action_status"] = string(constants.Pending)

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find pending cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}
