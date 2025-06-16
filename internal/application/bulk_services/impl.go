package bulkservices_application

import (
	"context"

	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

func (a *ApplicationStore) FetchServices(ctx context.Context, offset, limit int) ([]domain.Service, error) {
	data, err := a.service.GetServicePaginated(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (a *ApplicationStore) EnableDisableServicesMaker(ctx context.Context, Service_id string, action bool, makerId string) (string, error) {
	return a.service.UpdateServiceFlagRequest(ctx, Service_id, action, makerId)
}

func (a *ApplicationStore) EnableDisableServicesChecker(ctx context.Context, Action_id string, action bool, checkerId string) error {
	return a.service.UpdateServiceFlag(ctx, Action_id, action, checkerId)
}

func (a *ApplicationStore) FetchCifs(ctx context.Context, cif string) ([]domain.LinkedAccount, error) {
	return []domain.LinkedAccount{}, nil
}
func (a *ApplicationStore) RemoveCifMaker(ctx context.Context, id, makerId string) (string, error) {
	return "", nil
}
func (a *ApplicationStore) RemoveCifChecker(ctx context.Context, Action_id string, action bool, checkerId string) error {
	return nil
}
