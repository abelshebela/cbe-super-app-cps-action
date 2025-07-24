package accountlookup

import (
	"context"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
)

type UserSearchRepository interface {
	SearchUser(ctx context.Context, query string) (*domain.UserSearchResult, error)
}
