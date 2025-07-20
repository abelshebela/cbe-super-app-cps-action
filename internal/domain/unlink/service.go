// Package unlink provides the implementation for the service for unlink functionality.
package unlink

import (
	"context"

	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/unlink"
)

type Service struct {
	repository outbound.UnlinkRepository
}

func NewUnlinkService(repo outbound.UnlinkRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) UnlinkDevice(ctx context.Context, userCode string, cpsAction entities.CPSAction) (string, error) {
	action_code, err := s.repository.UnlinkDevice(ctx, userCode, cpsAction)
	if err != nil {
		return "", err
	}
	return action_code, nil
}

func (s *Service) Authorize(ctx context.Context, cpsAction *cps_entities.CPSAction) (*cps_entities.CPSAction, error) {
	return s.repository.Authorize(ctx, cpsAction)
}
