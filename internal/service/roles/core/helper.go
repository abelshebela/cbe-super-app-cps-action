package core

import (
	"cbe-super-app-cps-action/internal/storage"
	"context"
)

func CheckPortalCardsExistent(ctx context.Context, portalCards []string, portalCardRepo storage.PortalCardRepository) error {

	_, err := portalCardRepo.ValidatePortalCard(ctx, portalCards)
	if err != nil {
		return err
	}

	return nil
}
