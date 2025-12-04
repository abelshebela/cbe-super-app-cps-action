package merchant_lookup

import (
	"cbe-super-app-cps-action/internal/constants/dto/merchant_lookup"
	"cbe-super-app-cps-action/internal/constants/localization"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func ThreeClickMerchantLookup(ctx context.Context, client *http.Client, x_api_key, url, merchantId string) (*http.Response, merchant_lookup.MerchantLookUpResponse, error) {
	var accountInfo merchant_lookup.MerchantLookUpResponse
	var err error

	if merchantId == "" {
		return nil, merchant_lookup.MerchantLookUpResponse{}, localization.ErrorMerchantIDRequired
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+merchantId, nil)
	if err != nil {
		return nil, merchant_lookup.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", x_api_key)

	res, err := client.Do(req)
	if err != nil {
		return nil, merchant_lookup.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, merchant_lookup.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	if err := json.Unmarshal(resBody, &accountInfo); err != nil {
		return nil, merchant_lookup.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	return res, accountInfo, nil
}
