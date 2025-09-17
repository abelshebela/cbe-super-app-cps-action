package bank

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type BankPersistence interface {
	CreateBank(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	UpdateBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	DeleteBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	GetBank(ctx context.Context, id string) (*entity.Bank, error)
	GetAllBanks(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*entity.Bank], error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error)
	EnableOrDisableBank(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
	CheckExistingBank(ctx context.Context, name string) (bool, error)
	UpdateLogo(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
}
