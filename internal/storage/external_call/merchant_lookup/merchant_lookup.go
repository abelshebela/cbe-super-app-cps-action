package merchant_lookup

import (
	"cbe-super-app-cps-action/internal/constants/dto/merchant_lookup"
	"cbe-super-app-cps-action/internal/constants/localization"
	"context"
	"net/http"
	"time"

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
			Timeout: 30 * time.Second,
		},
	}
}

func (m *MerchantLookupAdapter) LookupMerchant(ctx context.Context, merchantID string) (merchant_lookup.MerchantLookUpResponse, error) {

	res, merchantInfo, err := ThreeClickMerchantLookup(ctx, m.client, m.x_api_key, m.Url, merchantID, m.logger)
	if err != nil {
		return merchant_lookup.MerchantLookUpResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return merchant_lookup.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	return merchantInfo, nil
}
