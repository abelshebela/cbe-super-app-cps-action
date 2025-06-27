package accountvalidation_app

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"
)

type ApplicationAbstracts interface {
	GetAccountValidation(ctx context.Context, id string) (account_validation.ValidationRule, error)
	UpdateAccountValidationRequest(ctx context.Context, id string, update account_validation.ValidationRule, makerID string, PhoneNumber string, FullName string) (string, error)
	UpdateAccountValidation(ctx context.Context, actionID string, action bool, checkerID string, PhoneNumber string, FullName string) error
}

type ApplicationStore struct {
	service account_validation.Service
	Logger  utils.Logger
}

func NewApplication(service account_validation.Service, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{service: service, Logger: logger}
}

func (a *ApplicationStore) GetAccountValidation(ctx context.Context, id string) (account_validation.ValidationRule, error) {

	return a.service.GetAccountValidation(ctx, id)
}

func (a *ApplicationStore) UpdateAccountValidationRequest(ctx context.Context, id string, update account_validation.ValidationRule, makerID string, PhoneNumber string, FullName string) (string, error) {
	return a.service.UpdateAccountValidationRequest(ctx, id, update, makerID, PhoneNumber, FullName)
}

func (a *ApplicationStore) UpdateAccountValidation(ctx context.Context, actionID string, action bool, checkerID string, PhoneNumber string, FullName string) error {
	return a.service.UpdateAccountValidation(ctx, actionID, action, checkerID, PhoneNumber, FullName)
}
