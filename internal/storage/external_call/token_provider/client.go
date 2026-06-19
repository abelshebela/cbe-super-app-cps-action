package token_provider

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TokenProviderClient struct {
	httpClient   *http.Client
	logger       utils.Logger
	tokenURL     string
	clientID     string
	clientSecret string
	scope        string
}

func NewTokenProviderClient(tokenURL, clientID, clientSecret, scope string, logger utils.Logger) *TokenProviderClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // UAT gateway uses internal CA not in system trust store
	}
	return &TokenProviderClient{
		httpClient:   &http.Client{Timeout: 30 * time.Second, Transport: transport},
		logger:       logger,
		tokenURL:     tokenURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		scope:        scope,
	}
}

func (c *TokenProviderClient) FetchToken(ctx context.Context) (*TokenResponse, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)
	form.Set("scope", c.scope)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Errorf("[TokenProvider] create request: %v", err)
		return nil, fmt.Errorf("token request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Errorf("[TokenProvider] http call: %v", err)
		return nil, fmt.Errorf("token fetch failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[TokenProvider] read body: %v", err)
		return nil, fmt.Errorf("token response read failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[TokenProvider] unexpected status %d: %s", resp.StatusCode, body)
		return nil, fmt.Errorf("token endpoint returned status %d", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		log.Errorf("[TokenProvider] unmarshal: %v", err)
		return nil, fmt.Errorf("token response parse failed: %w", err)
	}

	if tr.AccessToken == "" {
		return nil, fmt.Errorf("token endpoint returned empty access_token")
	}

	return &tr, nil
}
