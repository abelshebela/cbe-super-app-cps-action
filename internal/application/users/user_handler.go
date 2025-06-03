package users

import (
	"context"
	"cbe-super-app-member-users/internal/domain/users"
)

type ApplicationHandler interface {
	FetchLinkedAccounts(ctx context.Context, req FetchLinkedAccountsRequest) (*LinkedAccountResponseDto, error)
}

type applicationHandler struct {
	domainService users.UserService
}

func NewApplicationHandler(domainService users.UserService) ApplicationHandler {
	return &applicationHandler{domainService: domainService}
}

func (h *applicationHandler) FetchLinkedAccounts(ctx context.Context, req FetchLinkedAccountsRequest) (*LinkedAccountResponseDto, error) {
	response, err := h.domainService.ActiveLinkedAccounts(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return ToLinkedAccountResponseDto(response), nil
}