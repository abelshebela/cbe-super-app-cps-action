package amount_based_auth_app

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ApplicationService interface {
	UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth_domain.UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error)
	ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error)
	GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*amount_based_auth_domain.AuthTier], error)
}

type Handler struct {
	service *amount_based_auth_domain.Service
}

func AmountBasedAuthHandler(service *amount_based_auth_domain.Service) ApplicationService {
	return &Handler{
		service: service,
	}
}

func (h Handler) GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*amount_based_auth_domain.AuthTier], error) {
	authTiers, err := h.service.GetAllAmountBasedDetail(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return authTiers, nil
}
func (h Handler) UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth_domain.UpdateAmountBasedAuth,
	cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error) {
	return h.service.UpdateAmountBasedAuth(ctx, request, cpsAction)
}

func (h Handler) ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error) {
	return h.service.ApproveAmountBasedAuth(ctx, id, cpsAction)
}

func (h *Handler) RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error) {
	request, err := h.service.RejectAmountBasedAuth(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}
	return request, nil
}
