package accountlookup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountAPIClient struct {
	client     *http.Client
	logger     utils.Logger
	cbeBaseURL string
}

type AccountResponse struct {
	Data model.AccountInfo `json:"data"`
}
type Account interface {
	LookupAccountByPhone(ctx context.Context, phone string) (bool, error)
	LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountInfo, error)
}

// NewAccountAPIClient creates a new instance of AccountAPIClient
func InitAccountAPIClient(cbeBaseUrl string, timeout time.Duration, logger utils.Logger) Account {
	return &accountAPIClient{
		client:     &http.Client{Timeout: timeout},
		logger:     logger,
		cbeBaseURL: cbeBaseUrl,
	}
}

func (c *accountAPIClient) LookupAccountByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	accountReq := model.AccountLookUpRequest{
		PhoneNumber: phoneNumber,
	}

	data, err := json.Marshal(accountReq)
	if err != nil {
		c.logger.Errorf("failed to marshal account request %v", err)
		return false, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/%s", c.cbeBaseURL, "bps_banking/core/phone_number"), bytes.NewBuffer(data))
	if err != nil {
		c.logger.Errorf("failed to create HTTP request: %v", err)
		return false, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Errorf("HTTP request failed for phone number %s: %v", phoneNumber, err)
		return false, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		c.logger.Errorf("account lookup request failed with status code %v", res.StatusCode)
		return false, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: unexpected status code %d", res.StatusCode)
	}
}

func (b *accountAPIClient) LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountInfo, error) {

	data, err := json.Marshal(account)
	if err != nil {
		b.logger.Errorf("failed to marshal account request: %v", err)
		return nil, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/%s", b.cbeBaseURL, "bps_banking/core/account_number"), bytes.NewBuffer(data))
	if err != nil {
		b.logger.Errorf("failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := b.client.Do(req)
	if err != nil {
		if b.isTimeoutError(err, ctx) {
			b.logger.Errorf("request timeout for account %s: %v", account.AccountNumber, err)
			return nil, fmt.Errorf("TIME_OUT_ERROR: %w", err)
		}

		b.logger.Errorf("HTTP request failed for account %s: %v", account.AccountNumber, err)
		return nil, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b.logger.Errorf("request failed with status code %v", res.StatusCode)
		return nil, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: unexpected status code %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		b.logger.Errorf("failed to read response body for account %s: %v", account.AccountNumber, err)
		return nil, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}

	var accountInfo AccountResponse
	if err := json.Unmarshal(resBody, &accountInfo); err != nil {
		b.logger.Errorf("failed to unmarshal account info for account %s: %v", account.AccountNumber, err)
		return nil, fmt.Errorf("UNABLE_TO_CHECK_ACCOUNT: %w", err)
	}

	return &accountInfo.Data, nil
}

func (b *accountAPIClient) isTimeoutError(err error, ctx context.Context) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(ctx.Err(), context.DeadlineExceeded) ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "deadline exceeded")
}
