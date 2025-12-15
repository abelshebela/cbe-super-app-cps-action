package branch

import (
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type branchService struct {
	repo   storage.BranchRepository
	logger utils.Logger
}

func NewBranchService(repo storage.BranchRepository, logger utils.Logger) service.BranchService {
	return &branchService{
		repo:   repo,
		logger: logger,
	}
}

func (b *branchService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Branch", "Authorize")
	defer span.End()

	span.SetAttributes(attribute.String("action_code", cpsAction.ActionCode))

	b.logger.Infof("Branch service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	span.AddEvent("Branch action authorized", trace.WithAttributes(attribute.String("status", cpsAction.ActionStatus)))

	return cpsAction, nil
}
