package portalcard

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type portalCardDomain struct {
	portalInfra PortalCardRepository
}

func NewPortalCardDomain(poratal PortalCardRepository) PortaCardInterface {
	return &portalCardDomain{
		portalInfra: poratal,
	}
}

func (p *portalCardDomain) GetAllPortalCard(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*Card], error) {
	return p.portalInfra.GetAllPortalCard(ctx, filterParams)
}

func (p *portalCardDomain) ValidatePortalCard(ctx context.Context, names []string) (bool, error) {
	return p.portalInfra.ValidatePortalCard(ctx, names)
}
