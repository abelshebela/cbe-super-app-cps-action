// Package unlink provides the implementation for the service for unlink functionality.
package unlink

type Service struct {
	repository Repository
}

func NewUnlinkService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) UnlinkDevice(userCode, makerUser string, branchCode []string, homeBranch string) error {
	err := s.repository.UnlinkDevice(userCode, makerUser, branchCode, homeBranch)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ApproveOrDecline(userCode, decision, reason, checkerUser string) error {
	err := s.repository.ApproveOrDecline(userCode, decision, reason, checkerUser)
	if err != nil {
		return err
	}
	return nil
}
