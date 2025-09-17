package service

import (
	"context"

	portalcardDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PortalCardApplication interface {
	GetAll(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*portalcardDomain.Card], error)
}

type portalCardApp struct {
	portalDomain portalcardDomain.PortalCardRepository
	logger       utils.Logger
}

func NewPortalCardApp(service portalcardDomain.PortalCardRepository, logger utils.Logger) PortalCardApplication {
	return &portalCardApp{
		portalDomain: service,
		logger:       logger,
	}
}

func (s *portalCardApp) GetAll(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*portalcardDomain.Card], error) {
	return s.portalDomain.GetAllPortalCard(ctx, filterParams)
}
