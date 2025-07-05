package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

// AxiosSendSms sends an SMS message to the specified phone number using the configured notification service.
// Returns nil on success, or an error describing the failure.
func GetOtpExpiryTime(cfg config.VaultConfig) time.Time {
	// Set the OTP expiry time to 5 minutes from now
	wait, _ := strconv.Atoi(cfg.OtpWaitingTime)
	min := time.Duration(int64(wait))
	return time.Now().Add(min * time.Minute)
}
func AxiosSendSms(ctx context.Context, phoneNumber, messageBody string) error {
	ip, port, nodeEnv := os.Getenv("CONFIG_IP"), os.Getenv("CONFIG_PORT_LDAP_NOTIFICATION"), os.Getenv("GO_ENV")
	url := fmt.Sprintf("http://%s:%s/v1.0/chatbirrapi/ldapnotif/sms/send", ip, port)

	type requestBody struct {
		Recipient   string `json:"recipient"`
		MessageBody string `json:"messageBody"`
	}

	prefix := ""
	if nodeEnv == "uat" || nodeEnv == "dev" {
		prefix = "UAT: "
	}

	body := requestBody{
		Recipient:   FormatPhoneNumber(phoneNumber),
		MessageBody: prefix + messageBody,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling SMS body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating SMS request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS, status: %s", resp.Status)
	}
	return nil
}
