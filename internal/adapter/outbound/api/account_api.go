package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	// "os"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountAPIClient struct {
	// baseURL string
	client  *http.Client
	logger  utils.Logger
}

func NewAccountAPIClient(logger utils.Logger  ) *AccountAPIClient {

	// if baseURL == "" {
	// 	logger.Errorf("insert the account lookup url first")
	// 	return nil
	// }
	return &AccountAPIClient{
		// baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
		logger:  logger,
	}
}

func (c *AccountAPIClient) LookupAccountByPhone(ctx context.Context, phoneNumber string,PhoneLookupUrl string) (bool, error) {
	payload := map[string]string{"phone_number": phoneNumber}
	body, err := json.Marshal(payload)
	if err != nil {
		c.logger.Errorf("Failed to marshal phone lookup payload: %s", err.Error())
		return false, fmt.Errorf("API_REQUEST_FAILED")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, PhoneLookupUrl, bytes.NewBuffer(body))
	if err != nil {
		c.logger.Errorf("Failed to create phone lookup request: %s", err.Error())
		return false, fmt.Errorf("API_REQUEST_FAILED")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Errorf("Phone lookup API request failed: %s", err.Error())
		return false, fmt.Errorf("PHONE_LOOKUP_FAILED")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	c.logger.Errorf("Unexpected response status for phone lookup: %d", resp.StatusCode)
	return false, fmt.Errorf("API_REQUEST_FAILED")
}
