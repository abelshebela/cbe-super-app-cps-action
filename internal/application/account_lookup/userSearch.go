package accountlookup

import (
	"context"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
)

type UserSearchService struct {
	service domain.UserSearchRepository
}

func NewUserSearchService(service domain.UserSearchRepository) *UserSearchService {
	return &UserSearchService{service: service}
}

func (s *UserSearchService) SearchUser(ctx context.Context, query string) (*domain.UserSearchResult, error) {
	return s.service.SearchUser(ctx, query)
}
