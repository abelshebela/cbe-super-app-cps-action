package productcode

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constan "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

// Repository defines the interface for product code persistence operations
type Repository interface {
	FetchByID(ctx context.Context, id string) (*ProductCode, error)
	FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*ProductCode], error)
	Update(ctx context.Context, productCode *ProductCode) (*ProductCode, error)
}
