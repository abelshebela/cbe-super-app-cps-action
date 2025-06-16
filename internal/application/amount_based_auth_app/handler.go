package amount_based_auth_app

import amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/amount_based_auth"

type ApplicationService interface {
	UpdateAmountBasedAuth(request amount_based_auth_domain.AmountBasedAuthRequest) (amount_based_auth_domain.AuthTier, error)
}

type Handler struct {
	service *amount_based_auth_domain.Service
}

func (h Handler) UpdateAmountBasedAuth(request amount_based_auth_domain.AmountBasedAuthRequest) (amount_based_auth_domain.AuthTier, error) {
	tierRequest, err := h.service.BuildAuthTierRequest(request)
	if err != nil {
		return amount_based_auth_domain.AuthTier{}, err
	}

	return *tierRequest, nil
}

func AmountBasedAuthHandler(service *amount_based_auth_domain.Service) ApplicationService {
	return &Handler{
		service: service,
	}
}
