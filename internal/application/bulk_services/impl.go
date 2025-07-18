// Package bulkservicesApplication provides application logic for bulk services actions.
package bulkservicesApplication

import (
	"context"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationAbstracts interface {
	FetchServices(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, error)
	EnableDisableServicesMaker(ctx context.Context, ServiceID string, action bool, maker domain.User) (string, error)
	EnableDisableServicesChecker(ctx context.Context, ActionID string, action bool, checker domain.User) error
	FetchCifs(ctx context.Context, cif string) ([]domain.LinkedAccount, error)
	RemoveCifMaker(ctx context.Context, id []string, maker domain.User) (string, error)
	RemoveCifChecker(ctx context.Context, ActionID string, action bool, rejectionReason string, checker domain.User) (*domain.CPSAction, error)
}
type ApplicationStore struct {
	service domain.ServiceInterface
	Logger  utils.Logger
}

func NewAttachDetachChecker(service domain.ServiceInterface, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
		Logger:  logger,
	}
}

func (a *ApplicationStore) FetchServices(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, error) {
	data, err := a.service.GetServicePaginated(ctx, limit, offset)
	if err != nil {
		a.Logger.Errorf("[Application.FetchService] ", err.Error())
		return nil, err
	}
	return data, nil
}
func (a *ApplicationStore) EnableDisableServicesMaker(ctx context.Context, Service_id string, action bool, maker domain.User) (string, error) {
	data, err := a.service.UpdateServiceFlagRequest(ctx, Service_id, action, maker)
	if err != nil {
		a.Logger.Errorf("[Application.EnableDisableServicesMaker] ", err.Error())
		return "", err
	}
	return data, nil
}

func (a *ApplicationStore) EnableDisableServicesChecker(ctx context.Context, Action_id string, action bool, checker domain.User) error {
	err := a.service.UpdateServiceFlag(ctx, Action_id, action, checker)
	if err != nil {
		a.Logger.Errorf("[Application.EnableDisableServicesChecker] ", err.Error())

		return err
	}
	return nil
}

func (a *ApplicationStore) FetchCifs(ctx context.Context, cif string) ([]domain.LinkedAccount, error) {
	data, err := a.service.GetAccountByAccount(ctx, cif)
	if err != nil {
		a.Logger.Errorf("[Application.FetchCifs] ", err.Error())
		return nil, err
	}
	return data, nil
}
func (a *ApplicationStore) RemoveCifMaker(ctx context.Context, ids []string, maker domain.User) (string, error) {
	data, err := a.service.RemoveCifRequest(ctx, ids, true, maker)
	if err != nil {
		a.Logger.Errorf("[Application.RemoveCifMaker] ", err.Error())
		return "", err
	}
	return data, nil
}
func (a *ApplicationStore) RemoveCifChecker(ctx context.Context, Action_id string, action bool, rejectionReason string, checker domain.User) (*domain.CPSAction, error) {
	cpsAction, err := a.service.RemoveCif(ctx, Action_id, action, rejectionReason, checker)
	if err != nil {
		a.Logger.Errorf("[Application.RemoveCifChecker] ", err.Error())
		return nil, err
	}
	return cpsAction, nil
}
