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

	phoneURL := c.baseURL + "/bps_banking/core/fetch"
	body, err := json.Marshal(map[string]string{"number": number})
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ACCOUNT_NOT_FOUND")
	}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return &wrapper.Data, nil
}
