package service

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type ServiceRepo interface {
	GetAllService() ([]*ServiceResponse, error)
	GetOneService(id string) (*Service, error)
	UpdateOneService(id string, update UpdateServiceRequest) error
}

type ServiceDomain struct {
	repository Repository
}

func NewService(repository Repository) (ServiceRepo, error) {
	return &ServiceDomain{
		repository: repository,
	}, nil
}

func (s *ServiceDomain) GetAllService() ([]*ServiceResponse, error) {
	adverts, err := s.repository.GetAllService()
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (s *ServiceDomain) GetOneService(id string) (*Service, error) {
	advert, err := s.repository.GetOneService(id)
	if err != nil {
		return nil, err
	}

	return &advert, nil
}

func (s *ServiceDomain) UpdateOneService(id string, req UpdateServiceRequest) error {
	err := s.repository.UpdateOneService(id, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceDomain) CreateAction(ctx context.Context, req action.CPSAction) error {
	err := s.CreateAction(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
