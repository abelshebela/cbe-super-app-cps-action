package service

import (
	"context"

	portalcardDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PortalCardApplication interface {
	GetAll(ctx context.Context) ([]*portalcardDomain.Card, error)
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

func (s *portalCardApp) GetAll(ctx context.Context) ([]*portalcardDomain.Card, error) {
	service, err := s.portalDomain.GetAllPortalCard(ctx)
	if err != nil {
		return nil, err
	}
	return service, nil
}
