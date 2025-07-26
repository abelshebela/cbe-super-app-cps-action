package accountlookup

import "context"

type UserSearchService struct {
	repo UserSearchRepository
}

func NewUserSearchService(repo UserSearchRepository) *UserSearchService {
	return &UserSearchService{repo: repo}
}

func (s *UserSearchService) SearchUser(ctx context.Context, query string) (*UserSearchResult, error) {
	return s.repo.SearchUser(ctx, query)
}
