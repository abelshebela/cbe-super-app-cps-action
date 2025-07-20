package portalcard

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type portalCardDomain struct {
	portalInfra PortalCardRepository
}

func NewPortalCardDomain(poratal PortalCardRepository, logger utils.Logger) PortaCardInterface {
	return &portalCardDomain{
		portalInfra: poratal,
	}
}

func (p *portalCardDomain) GetAllPortalCard(ctx context.Context) ([]*Card, error) {

	portalcard, err := p.portalInfra.GetAllPortalCard(ctx)
	if err != nil {
		return nil, err
	}

	return portalcard, nil
}
