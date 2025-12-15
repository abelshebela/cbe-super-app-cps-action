package hq

import (
	"cbe-super-app-cps-action/internal/constants"
	hqDto "cbe-super-app-cps-action/internal/constants/dto/hq"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/hq/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type hqService struct {
	repo       storage.HQRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewHQService(repo storage.HQRepository, cpsActionService service.CPSActionService, logger utils.Logger) service.HQService {
	return &hqService{
		repo:       repo,
		cpsService: cpsActionService,
		logger:     logger,
	}
}

func (a *hqService) GetHQDetail(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.HQ], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetHQDetail", "HQ", "GetHQDetail")
	defer span.End()

	hqData, err := a.repo.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch HQ details", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return hqData, nil
}
func (a *hqService) GetHQ(ctx context.Context, id string) (hqDto.HQ, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetHQ", "HQ", "GetHQ")
	defer span.End()

	hq, err := a.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return hqDto.HQ{}, err
	}
	return hqDto.HQ{
		ID:          hq.ID.Hex(),
		BlockTime:   hq.BlockTime,
		ArchiveTime: hq.ArchiveTime,
	}, nil
}

func (a *hqService) GetBlockTime(ctx context.Context) (hqDto.BlockTimeResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetBlockTime", "HQ", "GetBlockTime")
	defer span.End()

	hq, err := a.repo.Find(ctx)
	if err != nil {
		a.logger.Errorf("failed to fetch HQ: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return hqDto.BlockTimeResponse{}, errors.New(localization.ErrorHQNotFound.Code)
	}
	return hqDto.BlockTimeResponse{
		BlockTime:      hq.BlockTime,
		CreatedAtBlock: hq.CreatedAtBlock,
		UpdatedAtBlock: hq.UpdatedAtBlock,
	}, nil
}

func (a *hqService) GetArchiveTime(ctx context.Context) (hqDto.ArchiveTimeResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetArchiveTime", "HQ", "GetArchiveTime")
	defer span.End()

	hq, err := a.repo.Find(ctx)
	if err != nil {
		a.logger.Errorf("failed to fetch HQ: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return hqDto.ArchiveTimeResponse{}, errors.New(localization.ErrorHQNotFound.Code)
	}
	return hqDto.ArchiveTimeResponse{
		ArchiveTime:      hq.ArchiveTime,
		CreatedAtArchive: hq.CreatedAtArchive,
		UpdatedAtArchive: hq.UpdatedAtArchive,
	}, nil

}

func (a *hqService) GetPasswordExpiry(ctx context.Context) (hqDto.PasswordExpiryResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPasswordExpiry", "HQ", "GetPasswordExpiry")
	defer span.End()

	hq, err := a.repo.Find(ctx)
	if err != nil {
		a.logger.Errorf("failed to fetch HQ: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return hqDto.PasswordExpiryResponse{}, errors.New(localization.ErrorHQNotFound.Code)
	}
	return hqDto.PasswordExpiryResponse{
		PasswordExpiry:          hq.PasswordExpiry,
		CreatedAtPasswordExpiry: hq.CreatedAtPasswordExpiry,
		UpdatedAtPasswordExpiry: hq.UpdatedAtPasswordExpiry,
	}, nil
}

func (a *hqService) UpdateBlockTime(ctx context.Context, request hqDto.UpdateBlockTimeRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBlockTime", "HQ", "UpdateBlockTime")
	defer span.End()

	originalHQ, err := a.repo.Find(ctx)
	if err != nil {
		a.logger.Errorf("failed to fetch HQ: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return errors.New(localization.ErrorHQNotFound.Code)
	}

	updatedHQ := *originalHQ
	updatedHQ.BlockTime = request.BlockTime

	if err := core.HandleCPSAction(ctx, a.cpsService, originalHQ.ID.Hex(), constants.RequestUpdateHQBlockTime, updatedHQ, *originalHQ, constants.ActionUpdate); err != nil {
		a.logger.Errorf("CPS action failed for BlockTime %s: %v", updatedHQ.BlockTime, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", originalHQ.ID.Hex()),
		))
		return err
	}
	return nil
}

func (a *hqService) UpdateArchiveTime(ctx context.Context, request hqDto.UpdateArchiveTimeRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateArchiveTime", "HQ", "UpdateArchiveTime")
	defer span.End()

	originalHQ, err := a.repo.Find(ctx)
	if err != nil {
		a.logger.Errorf("failed to fetch HQ: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return errors.New(localization.ErrorHQNotFound.Code)
	}

	updatedHQ := *originalHQ
	updatedHQ.ArchiveTime = request.ArchiveTime

	if err := core.HandleCPSAction(ctx, a.cpsService, originalHQ.ID.Hex(), constants.RequestUpdateHQArchiveTime, updatedHQ, *originalHQ, constants.ActionUpdate); err != nil {
		a.logger.Errorf("CPS action failed for ArchiveTime %v: %v", updatedHQ.ArchiveTime, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", originalHQ.ID.Hex()),
		))
		return err
	}

	return nil
}

func (a *hqService) UpdatePasswordExpiry(ctx context.Context, request hqDto.UpdatePasswordExpiryRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdatePasswordExpiry", "HQ", "UpdatePasswordExpiry")
	defer span.End()

	originalHQ, err := a.repo.Find(ctx)
	if err != nil {
		a.logger.Errorf("failed to fetch HQ: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return errors.New(localization.ErrorHQNotFound.Code)
	}

	updatedHQ := *originalHQ
	updatedHQ.PasswordExpiry = request.PasswordExpiry

	if err := core.HandleCPSAction(ctx, a.cpsService, originalHQ.ID.Hex(), constants.RequestUpdatePasswordExpiry, updatedHQ, *originalHQ, constants.ActionUpdate); err != nil {
		a.logger.Errorf("CPS action failed for PasswordExpiry %v: %v", updatedHQ.PasswordExpiry, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", originalHQ.ID.Hex()),
		))
		return err
	}

	return nil
}

func (s *hqService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "HQ", "Authorize")
	defer span.End()

	requestedAction := action.RequestAction

	hq, err := local_util.JsonUnmarshal[model.HQ](action.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch requestedAction {
	case string(constants.RequestUpdateHQBlockTime):
		err = s.repo.Update(ctx, "block_time", hq.BlockTime, time.Now())
		if err != nil {
			span.AddEvent("Failed to update block time", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateHQArchiveTime):
		err = s.repo.Update(ctx, "archive_time", hq.ArchiveTime, time.Now())
		if err != nil {
			span.AddEvent("Failed to update archive time", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdatePasswordExpiry):
		err = s.repo.Update(ctx, "password_expiry", hq.PasswordExpiry, time.Now())
		if err != nil {
			span.AddEvent("Failed to update password expiry", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	default:
		span.AddEvent("Invalid request", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(requestedAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	action.CurrentAction = hq
	action.ActionStatus = constants.Approved
	s.logger.Infof("HQ action approved", "action_code", action.ActionCode, "request", requestedAction)
	return action, nil
}
