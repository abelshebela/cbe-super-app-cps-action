package unlink

import "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/unlink/entities"

type Service struct {
	repository Repository
}

func NewUnlinkService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) UnlinkDevice(userCode string, cpsAction entities.CPSAction) error {
	err := s.repository.UnlinkDevice(userCode, cpsAction)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ApproveOrDecline(userCode, decision, reason string, cpsAction entities.CPSAction) error {
	err := s.repository.ApproveOrDecline(userCode, decision, reason, cpsAction)
	if err != nil {
		return err
	}
	return nil
}
