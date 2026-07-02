package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	bpsactionsvc "cbe-super-app-cps-action/internal/service/bps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"slices"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/trace"
)

func ParameterProvider(ctx context.Context, filterParams *types.Filter, span trace.Span, role string, idxRepo storage.BPSActionApproveIndexRepository, log utils.Logger) (*types.Filter, []string, error) {
	var actionData []string
	rawRoleID, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		return nil, nil, errors.New(localization.ErrorOperationNotAllowed.Message)
	}

	_, makerActions, checkerActions, auditorActions, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		return nil, nil, errors.New(err.Error())
	}

	if role == constants.Maker {
		actionData = makerActions
	} else if role == constants.Checker {
		actionData = checkerActions
	} else if role == constants.Auditor {
		actionData = auditorActions
	}

	if actionData == nil {
		return nil, nil, errors.New(localization.ErrorOperationNotAllowed.Message)
	}

	// resolve action_names -> request_actions
	var filterReqs, reqs []string
	seen := map[string]struct{}{}

	services := local_util.ExtractStringSlice(filterParams.Filters, "services")
	log.Infof("[CPSAction][GetUserApproverActions] list of service trying to filer ******** services:%s", services)

	log.Infof("[CpsActionH][Approve] checker actions: %v", actionData)

	if len(services) != 0 {
		for _, chk := range services {
			if !slices.Contains(actionData, chk) {
				return nil, nil, errors.New(localization.ErrorNotAllowedServicesIncluded.Message)
			}
		}

		for _, svc := range services {
			if lst, ok := bpsactionsvc.RequestActionGroups[strings.ToUpper(svc)]; ok {
				for _, ra := range lst {
					key := string(ra)
					if _, ok := seen[key]; ok {
						continue
					}

					seen[key] = struct{}{}
					filterReqs = append(filterReqs, key)
				}
			}
		}
	}

	seen = map[string]struct{}{}
	for _, mod := range actionData {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
			for _, ra := range lst {
				key := string(ra)
				if _, ok := seen[key]; ok {
					continue
				}
				if len(services) > 0 {
					if !slices.Contains(filterReqs, key) {
						continue
					}
				}

				seen[key] = struct{}{}
				reqs = append(reqs, key)
			}
		}
	}

	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	if len(reqs) > 0 {
		filterParams.Filters["request_action"] = map[string]interface{}{"$in": reqs}
	}

	return filterParams, reqs, nil
}
