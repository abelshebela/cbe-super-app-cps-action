package amount_based_auth_app

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
)

type ApplicationService interface {
	UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth_domain.UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectCPSAction) (*model.CPSAction, error)
}

type Handler struct {
	service *amount_based_auth_domain.Service
}

func AmountBasedAuthHandler(service *amount_based_auth_domain.Service) ApplicationService {
	return &Handler{
		service: service,
	}
}

func (h Handler) UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth_domain.UpdateAmountBasedAuth,
	cpsAction model.CreateCPSAction) (*model.CPSAction, error) {
	cpsActionRes, err := h.service.UpdateAuthTier(ctx, request, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (h Handler) ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error) {
	request, err := h.service.ApproveAuthTierApprove(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (h *Handler) RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectCPSAction) (*model.CPSAction, error) {
	request, err := h.service.RejectAuthTier(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}
	return request, nil
}
