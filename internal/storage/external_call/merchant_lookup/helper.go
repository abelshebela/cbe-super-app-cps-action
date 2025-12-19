package merchant_lookup

import (
	merchantDto "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ThreeClickMerchantLookup(ctx context.Context, client *http.Client, x_api_key, token, url, merchantId string, logger utils.Logger) (*http.Response, merchantDto.MerchantLookUpResponse, error) {
	var apiResp merchantDto.MerchantLookupAPIResponse

	if merchantId == "" {
		return nil, merchantDto.MerchantLookUpResponse{}, localization.ErrorMerchantIDRequired
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+merchantId, nil)
	if err != nil {
		logger.Errorf("error with context error: %v ", err)
		return nil, merchantDto.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", x_api_key)
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error while requesting merchant lookup error :%v", err)
		return nil, merchantDto.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Errorf("error reading body: %v", err)
		return nil, merchantDto.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		logger.Errorf("system error while unmarshaling error:%v", err)
		return nil, merchantDto.MerchantLookUpResponse{}, localization.ErrorUnexpectedError
	}

	if len(apiResp.Data) == 0 {
		return res, merchantDto.MerchantLookUpResponse{}, fmt.Errorf("no merchant found")
	}

	merchant := apiResp.Data[0]
	return res, merchant, nil
}
