package bpscalls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BpsCallsInterface interface {
	FetchLinkedAccount(accountNumber string) (*model.LinkedAccount, error)
}
type BpsCalls struct {
	cfg    *config.VaultConfig
	logger utils.Logger
}

func NewBpsCalls(cfg *config.VaultConfig, logger utils.Logger) BpsCallsInterface {
	return &BpsCalls{
		cfg:    cfg,
		logger: logger,
	}
}

func (b *BpsCalls) FetchLinkedAccount(accountNumber string) (*model.LinkedAccount, error) {
	// url := b.cfg.BPSAccountFetchURL //return this line and remove second
	url := b.cfg.AccountLookUpURI

	if url == "" {
		b.logger.Errorf("missing BPS_ACCOUNT_FETCH_URL in config")
		return nil, fmt.Errorf("MISSING_CONFIG_KEY")
	}

	b.logger.Infof("fetching linked account", "url", url, "account_number", accountNumber)

	// Create payload
	payload := strings.NewReader(fmt.Sprintf(`{"account_number": "%s"}`, accountNumber))

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		b.logger.Errorf("failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Execute request
	res, err := client.Do(req)
	if err != nil {
		b.logger.Errorf("HTTP request failed", "error", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
	defer res.Body.Close()

	// Check HTTP status code
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			b.logger.Warnf("linked account not found", "account_number", accountNumber)
			return nil, fmt.Errorf("ACCOUNT_NOT_FOUND")
		}
		b.logger.Warnf("unexpected HTTP status code", "status_code", res.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Check Content-Type header
	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		b.logger.Warnf("unexpected Content-Type in response", "content_type", contentType)
		return nil, fmt.Errorf("unexpected Content-Type: %s", contentType)
	}

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		b.logger.Errorf("failed to read response body", "error", err)
		return nil, err
	}

	// Validate JSON
	if !json.Valid(body) {
		b.logger.Errorf("response is not valid JSON", "body", string(body))
		return nil, fmt.Errorf("invalid JSON response")
	}

	// Unmarshal response
	var linkedAccountExt model.LinkedAccountExternal
	if err := json.Unmarshal(body, &linkedAccountExt); err != nil {
		b.logger.Errorf("failed to unmarshal linked account response", "error", err, "body", string(body))
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	b.logger.Infof("linked account fetched successfully")

	linkedAccount := linkedAccountExt.MapFromExternal()
	return &linkedAccount, nil
}
