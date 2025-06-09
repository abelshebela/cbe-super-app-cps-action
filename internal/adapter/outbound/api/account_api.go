package api

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "os"
    "time"

    // domainAccount "cbe-super-app-member-users/internal/domain/account"
    // accountPort "cbe-super-app-member-users/port/outbound/account"
    local "cbe-super-app-member-users/internal/shared"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountAPIClient struct {
    baseURL string
    client  *http.Client
    logger  utils.Logger
}

func NewAccountAPIClient(logger utils.Logger) *AccountAPIClient {
    baseURL := os.Getenv("API_BASE_URL")
    if baseURL == "" {
        baseURL = "https://devcbe.eaglelionsystems.com/api/v1/cbesuperapp/bps/core"
    }
    return &AccountAPIClient{
        baseURL: baseURL,
        client:  &http.Client{Timeout: 30 * time.Second},
        logger:  logger,
    }
}

func (c *AccountAPIClient) LookupAccountByPhone(ctx context.Context, phoneNumber string) (bool, error) {
    payload := map[string]string{"phone_number": phoneNumber}
    body, err := json.Marshal(payload)
    if err != nil {
        c.logger.Errorf("Failed to marshal phone lookup payload: %s", err.Error())
        return false, local.DefineError.Account["API_REQUEST_FAILED"]
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/account_number", bytes.NewBuffer(body))
    if err != nil {
        c.logger.Errorf("Failed to create phone lookup request: %s", err.Error())
        return false, local.DefineError.Account["API_REQUEST_FAILED"]
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        c.logger.Errorf("Phone lookup API request failed: %s", err.Error())
        return false, local.DefineError.Account["PHONE_LOOKUP_FAILED"]
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        return true, nil // Account exists
    } else if resp.StatusCode == http.StatusNotFound {
        return false, nil // Account doesn't exist
    }

    c.logger.Errorf ("Unexpected response status for phone lookup: %d", resp.StatusCode)
    return false, local.DefineError.Account["API_REQUEST_FAILED"]
}
