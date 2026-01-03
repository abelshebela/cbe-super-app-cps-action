package merchant_lookup

import (
	"bytes"
	merchantDto "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ThreeClickMerchantLookup(ctx context.Context, client *http.Client, x_api_key, url, merchantId string, logger utils.Logger) (*http.Response, merchantDto.MerchantLookUpResponse, error) {
	var apiResp merchantDto.MerchantLookupAPIResponse

	if merchantId == "" {
		return nil, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorMerchantIDRequired.Code)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+merchantId, nil)
	if err != nil {
		logger.Errorf("error with context error: %v ", err)
		return nil, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", x_api_key)

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error while requesting merchant lookup error :%v", err)
		return nil, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Errorf("error reading body: %v", err)
		return nil, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		logger.Errorf("system error while unmarshaling error:%v", err)
		return nil, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if apiResp.Status == 401 {
		return res, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorThirdAPIRequestNotUnauthorized.Code)
	}

	if len(apiResp.Data) == 0 {
		return res, merchantDto.MerchantLookUpResponse{}, errors.New(localization.ErrorResourceNotFound.Code)
	}

	merchant := apiResp.Data[0]
	return res, merchant, nil
}

func ThreeClickMerchantUpdate(ctx context.Context, client *http.Client, x_api_key, url, merchantId string, payload merchantDto.ERPUpdateMerchantRequest, logger utils.Logger) (*http.Response, error) {
	if merchantId == "" {
		return nil, errors.New(localization.ErrorMerchantIDRequired.Code)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("failed to marshal ERP update payload: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url+merchantId, bytes.NewReader(body))
	if err != nil {
		logger.Errorf("error with context error: %v ", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", x_api_key)

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error while requesting merchant update error :%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return res, nil
}
