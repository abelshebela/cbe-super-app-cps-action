package account_lookup

import (
	"bytes"
	accountLookup "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func BPSBankingClient(ctx context.Context, client *http.Client, inputType string, inputData any, url string) (*http.Response, accountLookup.AccountResponse, error) {
	var accountInfo accountLookup.AccountResponse
	var data []byte
	var err error
	if inputType == "phone_number" {
		if _, ok := inputData.(string); !ok {
			return nil, accountLookup.AccountResponse{}, errors.New("input data is not string")
		}
		input := accountLookup.AccountLookUpRequest{
			PhoneNumber: inputData.(string),
		}
		data, err = json.Marshal(input)
		if err != nil {
			return nil, accountLookup.AccountResponse{}, err
		}
	} else if inputType == "account_number" {
		if _, ok := inputData.(string); !ok {
			return nil, accountLookup.AccountResponse{}, errors.New("input data is not string")
		}
		input := accountLookup.AccountLookUpRequest{
			AccountNumber: inputData.(string),
		}
		data, err = json.Marshal(input)
		if err != nil {
			return nil, accountLookup.AccountResponse{}, err
		}
	} else if inputType == "with_fayda" {
		if _, ok := inputData.(accountLookup.CreateAccountRequest); !ok {
			return nil, accountLookup.AccountResponse{}, errors.New("input data is not CreateAccountRequest")
		}
		input := inputData.(accountLookup.CreateAccountRequest)
		data, err = json.Marshal(input)
		if err != nil {
			return nil, accountLookup.AccountResponse{}, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, accountLookup.AccountResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, accountLookup.AccountResponse{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, accountLookup.AccountResponse{}, err
	}

	if err := json.Unmarshal(resBody, &accountInfo); err != nil {
		return nil, accountLookup.AccountResponse{}, err
	}

	return res, accountInfo, nil
}
