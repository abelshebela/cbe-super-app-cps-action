package donation

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type donationService struct {
	repo   storage.DonationRepository
	logger utils.Logger
}

func NewDonationService(repo storage.DonationRepository, logger utils.Logger) service.DonationService {
	return &donationService{
		repo:   repo,
		logger: logger,
	}
}

func (d *donationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	d.logger.Infof("Donation service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
