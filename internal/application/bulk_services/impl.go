package bulkservices_application

import (
	"context"

	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

func (a *ApplicationStore) FetchServices(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, error) {
	data, err := a.service.GetServicePaginated(ctx, limit, offset)
	if err != nil {
		a.Logger.Errorf("[Application.FetchService] ", err.Error())
		return nil, err
	}
	return data, nil
}
func (a *ApplicationStore) EnableDisableServicesMaker(ctx context.Context, Service_id string, action bool, makerId string) (string, error) {
	data, err := a.service.UpdateServiceFlagRequest(ctx, Service_id, action, makerId)
	if err != nil {
		a.Logger.Errorf("[Application.EnableDisableServicesMaker] ", err.Error())
		return "", err
	}
	return data, nil
}

func (a *ApplicationStore) EnableDisableServicesChecker(ctx context.Context, Action_id string, action bool, checkerId string) error {
	err := a.service.UpdateServiceFlag(ctx, Action_id, action, checkerId)
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
func (a *ApplicationStore) RemoveCifMaker(ctx context.Context, id []string, makerId string) (string, error) {
	data, err := a.service.RemoveCifRequest(ctx, id, true, makerId)
	if err != nil {
		a.Logger.Errorf("[Application.RemoveCifMaker] ", err.Error())
		return "", err
	}
	return data, nil
}
func (a *ApplicationStore) RemoveCifChecker(ctx context.Context, Action_id string, action bool, checkerId string) error {
	err := a.service.RemoveCif(ctx, Action_id, action, checkerId)
	if err != nil {
		a.Logger.Errorf("[Application.RemoveCifChecker] ", err.Error())
		return err
	}
	return nil
}