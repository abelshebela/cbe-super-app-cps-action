package amount_based_auth_app

import (
	amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/amount_based_auth"
)

type ApplicationService interface {
	UpdateAmountBasedAuth(request amount_based_auth_domain.AmountBasedAuthRequest) (string, error)
	ApproveAmountBasedAuth(id string) (string, error)
}

type Handler struct {
	service *amount_based_auth_domain.Service
}

func (h Handler) ApproveAmountBasedAuth(id string) (string, error) {
	request, err := h.service.BuildAuthTierApproveRequest(id)
	if err != nil {
		return "", err
	}

	return request, nil
}

func (h Handler) UpdateAmountBasedAuth(request amount_based_auth_domain.AmountBasedAuthRequest) (string, error) {
	tierRequest, err := h.service.BuildAuthTierRequest(request)
	if err != nil {
		return "Something went wrong", err
	}

	return *tierRequest, nil
}

func AmountBasedAuthHandler(service *amount_based_auth_domain.Service) ApplicationService {
	return &Handler{
		service: service,
	}
}
