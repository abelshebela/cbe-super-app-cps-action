package bpscalls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_lookup"
)

type CBEUserSearchClient struct {
	baseURL string
	client  *http.Client
}

func NewCBEUserSearchClient(baseURL string) inbound.UserSearchRepository {
	return &CBEUserSearchClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *CBEUserSearchClient) SearchUser(ctx context.Context, number string) (*domain.UserSearchResult, error) {
	type userSearchResponse struct {
		Data domain.UserSearchResult `json:"data"`
	}

	var wrapper userSearchResponse

	// 1st attempt: Search by phone number
	phoneURL := c.baseURL + "/bps_banking/core/phone_number"
	body, err := json.Marshal(map[string]string{"phone_number": number})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, phoneURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
			return nil, err
		}
		return &wrapper.Data, nil
	}

	// 2nd attempt: fallback search by account number
	accountURL := c.baseURL + "/bps_banking/core/account_number"
	body, err = json.Marshal(map[string]string{"account_number": number})
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, accountURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ACCOUNT_NOT_FOUND")
	}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return &wrapper.Data, nil
}
