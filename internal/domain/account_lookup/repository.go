package accountlookup

import (
	"context"
)

type UserSearchRepository interface {
	SearchUser(ctx context.Context, query string) (*UserSearchResult, error)
}
