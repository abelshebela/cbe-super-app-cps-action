package users

import (
    "context"
    "cbe-super-app-member-users/internal/domain/users"
)

type ApplicationService interface {
    FetchLinkedAccounts(ctx context.Context, id string) (*users.LinkedAccountResponse, error)
}

type applicationService struct {
    domainService *users.UserService
}

func NewApplicationService(domainService *users.UserService) ApplicationService {
    return &applicationService{domainService: domainService}
}

func (s *applicationService) FetchLinkedAccounts(ctx context.Context, id string) (*users.LinkedAccountResponse, error) {
    return s.domainService.ActiveLinkedAccounts(ctx, id)
}