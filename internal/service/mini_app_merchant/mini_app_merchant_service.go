package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type miniAppMerchantService struct {
	repo   storage.MiniAppMerchantRepository
	logger utils.Logger
}

func NewMiniAppMerchantService(repo storage.MiniAppMerchantRepository, logger utils.Logger) service.MiniAppMerchantService {
	return &miniAppMerchantService{
		repo:   repo,
		logger: logger,
	}
}

func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	m.logger.Infof("Mini App Merchant service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

func (m *miniAppMerchantService) DetailMiniAppByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {

	m.logger.Infof("Mini App Merchant service authorizing action: %s", id)

	return m.repo.DetailMiniAppByID(ctx, id)
}

func (s *miniAppMerchantService) AddMiniApp(ctx context.Context, merchantID string, miniApp model.MiniApps) error {
	return s.repo.AddMiniApp(ctx, merchantID, miniApp)
}
func (s *miniAppMerchantService) UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error {
	return s.repo.UpdateMiniAppEnabledState(ctx, merchantID, miniAppID, enabled)
}
func (s *miniAppMerchantService) SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error {
	return s.repo.SoftDeleteMiniApp(ctx, merchantID, miniAppID)
}
