package accountvalidation_app

import (
	"context"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type ApplicationAbstracts interface {
	GetAccountValidation(ctx context.Context, id string) (dto.GetAccountValidationResponse, error)
	UpdateAccountValidationRequest(ctx context.Context, id string, update account_validation.ValidationRule, maker account_validation.User) (dto.UpdateAccountValidationResponse, error)
	UpdateAccountValidation(ctx context.Context, actionID string, decision utils.DecisonEnum, checker account_validation.User, rejectedReason string) error
}

type ApplicationStore struct {
	service account_validation.Service
	Logger  sharedutils.Logger
}

func NewApplication(service account_validation.Service, logger sharedutils.Logger) ApplicationAbstracts {
	return &ApplicationStore{service: service, Logger: logger}
}

func (a *ApplicationStore) GetAccountValidation(ctx context.Context, id string) (dto.GetAccountValidationResponse, error) {
	rule, err := a.service.GetAccountValidation(ctx, id)
	if err != nil {
		return dto.GetAccountValidationResponse{}, err
	}
	return dto.GetAccountValidationResponse{Validation: dto.ToValidationRuleDTO(rule)}, nil
}

func (a *ApplicationStore) UpdateAccountValidationRequest(ctx context.Context, id string, update account_validation.ValidationRule, maker account_validation.User) (dto.UpdateAccountValidationResponse, error) {
	actionID, err := a.service.UpdateAccountValidationRequest(ctx, id, update, maker)
	if err != nil {
		return dto.UpdateAccountValidationResponse{}, err
	}
	return dto.UpdateAccountValidationResponse{ActionID: actionID}, nil
}

func (a *ApplicationStore) UpdateAccountValidation(ctx context.Context, actionID string, decision utils.DecisonEnum, checker account_validation.User, rejectedReason string) error {
	return a.service.UpdateAccountValidation(ctx, actionID, decision, checker, rejectedReason)
}
