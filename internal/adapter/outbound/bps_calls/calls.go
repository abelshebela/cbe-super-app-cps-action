package bpscalls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type BpsCallsInterface interface {
	FetchLinkedAccount(accountNumber string) (*model.LinkedAccount, error)
}
type BpsCalls struct{}

func NewBpsCalls() BpsCallsInterface {
	return &BpsCalls{}
}

func (b *BpsCalls) FetchLinkedAccount(accountNumber string) (*model.LinkedAccount, error) {
	url := "https://devcbe.eaglelionsystems.com/api/v1/cbesuperapp/bps/core/account_number"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
		"account_number": "%s"
	}`, accountNumber))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var linkedAccount model.LinkedAccount
	if err := json.Unmarshal(body, &linkedAccount); err != nil {
		return nil, err
	}

	return &linkedAccount, nil
}
