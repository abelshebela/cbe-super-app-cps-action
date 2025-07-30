package productcode

import (
	"context"

	domian "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constan "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

)

// Repository defines the interface for product code persistence operations
type Repository interface {
	FetchByID(ctx context.Context, id string) (*domian.ProductCode, error)
	FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*domian.ProductCode], error)
	Update(ctx context.Context, productCode *domian.ProductCode) (*domian.ProductCode, error)
}
