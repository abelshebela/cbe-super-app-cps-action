package merchant_lookup

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"context"
	"net/http"
	"time"

	merchantDto "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MerchantLookupAdapter struct {
	Url       string
	cfg       config.VaultConfig
	x_api_key string
	logger    utils.Logger
	client    *http.Client
}

func NewMerchantLookupAdapter(url string, cfg config.VaultConfig, x_api_key string, logger utils.Logger) *MerchantLookupAdapter {
	return &MerchantLookupAdapter{
		Url:       url,
		cfg:       cfg,
		logger:    logger,
		x_api_key: x_api_key,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (m *MerchantLookupAdapter) LookupMerchant(ctx context.Context, merchantID, token string) (merchantDto.MerchantLookUpResponse, error) {

	res, merchantInfo, err := ThreeClickMerchantLookup(ctx, m.client, m.x_api_key, token, m.Url, merchantID, m.logger)
	if err != nil {
		m.logger.Errorf("failed to lookup merchant data from third party API: %v", err)
		return merchantDto.MerchantLookUpResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return merchantDto.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	return merchantInfo, nil
}


