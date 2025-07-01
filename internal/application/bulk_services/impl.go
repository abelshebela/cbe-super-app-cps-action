package bulkservices_application

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"

	domain "cbe-super-app-cps-action/internal/domain/action"
)

func (a *ApplicationStore) FetchServices(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, *common.ErrorDefinition) {
	data, err := a.service.GetServicePaginated(ctx, limit, offset)
	if err != nil {
		err_def := common.ErrorDefinition{
			Code:    "",
			Message: "",
		}
		a.Logger.Errorf("[Application.FetchService] ", err.Error())
		a.Logger.Errorf("[Application.FetchService] ", err_def)
		return nil, &err_def
	}
	return data, nil
}
func (a *ApplicationStore) EnableDisableServicesMaker(ctx context.Context, Service_id string, action bool, makerId string) (string, *common.ErrorDefinition) {
	data, err := a.service.UpdateServiceFlagRequest(ctx, Service_id, action, makerId)
	if err != nil {
		err_def := common.ErrorDefinition{
			Code:    "",
			Message: "",
		}
		a.Logger.Errorf("[Application.EnableDisableServicesMaker] ", err.Error())
		a.Logger.Errorf("[Application.EnableDisableServicesMaker] ", err_def)
		return "", &err_def
	}
	return data, nil
}

func (a *ApplicationStore) EnableDisableServicesChecker(ctx context.Context, Action_id string, action bool, checkerId string) *common.ErrorDefinition {
	err := a.service.UpdateServiceFlag(ctx, Action_id, action, checkerId)
	if err != nil {
		err_def := common.ErrorDefinition{
			Code:    "",
			Message: "",
		}
		a.Logger.Errorf("[Application.EnableDisableServicesChecker] ", err.Error())
		a.Logger.Errorf("[Application.EnableDisableServicesChecker] ", err_def)
		return &err_def
	}
	return nil
}

func (a *ApplicationStore) FetchCifs(ctx context.Context, cif string) ([]domain.LinkedAccount, *common.ErrorDefinition) {
	data, err := a.service.GetAccountByAccount(ctx, cif)
	if err != nil {
		err_def := common.ErrorDefinition{
			Code:    "",
			Message: "",
		}
		a.Logger.Errorf("[Application.FetchCifs] ", err.Error())
		a.Logger.Errorf("[Application.FetchCifs] ", err_def)
		return nil, &err_def
	}
	return data, nil
}
func (a *ApplicationStore) RemoveCifMaker(ctx context.Context, id []string, makerId string) (string, *common.ErrorDefinition) {
	data, err := a.service.RemoveCifRequest(ctx, id, true, makerId)
	if err != nil {
		err_def := common.ErrorDefinition{
			Code:    "",
			Message: "",
		}
		a.Logger.Errorf("[Application.RemoveCifMaker] ", err.Error())
		a.Logger.Errorf("[Application.RemoveCifMaker] ", err_def)
		return "", &err_def
	}
	return data, nil
}
func (a *ApplicationStore) RemoveCifChecker(ctx context.Context, Action_id string, action bool, checkerId string) *common.ErrorDefinition {
	err := a.service.RemoveCif(ctx, Action_id, action, checkerId)
	if err != nil {
		err_def := common.ErrorDefinition{
			Code:    "",
			Message: "",
		}
		a.Logger.Errorf("[Application.RemoveCifChecker] ", err.Error())
		a.Logger.Errorf("[Application.RemoveCifChecker] ", err_def)
		return &err_def
	}
	return nil
}
