package merchant_lookup

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"context"
	"errors"
	"net/http"
	"time"

	merchantDto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"

	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MerchantLookupAdapter struct {
	cfg    config.VaultConfig
	logger utils.Logger
	client *http.Client
}

func NewMerchantLookupAdapter(cfg config.VaultConfig, logger utils.Logger) *MerchantLookupAdapter {
	return &MerchantLookupAdapter{
		cfg:    cfg,
		logger: logger,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (m *MerchantLookupAdapter) LookupMerchant(ctx context.Context, merchantID, xAPIKey, url string) (merchantDto.MerchantLookUpResponse, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	res, merchantInfo, err := ThreeClickMerchantLookup(ctx, m.client, xAPIKey, url, merchantID, m.logger)
	if err != nil {
		log.Errorf("failed to lookup merchant data from third party API: %v", err)
		return merchantDto.MerchantLookUpResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return merchantInfo, nil
}

func (m *MerchantLookupAdapter) UpdateMerchant(ctx context.Context, merchantID string, payload merchantDto.ERPUpdateMerchantRequest, xAPIKey, url string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	res, err := ThreeClickMerchantUpdate(ctx, m.client, xAPIKey, url, merchantID, payload, m.logger)
	if err != nil {
		log.Errorf("failed to update merchant data to third party API: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Errorf("unexpected status code from third party API: %d", res.StatusCode)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
