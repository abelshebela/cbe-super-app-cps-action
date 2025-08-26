package external_call

import (
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SMSPersistence struct {
	httpClient *http.Client
	logger     utils.Logger
	baseURL    string
}

// // baseURL : "https://devcbe.eaglelionsystems.com/api/v1.0/chatbirrapi/ldapnotif/sms/send"
func NewSMSPersistence(baseUrl string, logger utils.Logger) *SMSPersistence {
	return &SMSPersistence{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:  logger,
		baseURL: baseUrl,
	}
}

// // SendSMS sends an SMS using the external API
// func (s *SMSPersistence) SendSMS(ctx context.Context, recipient, messageBody string) error {

// 	payload := dto.SMSRequest{
// 		Recipient:   recipient,
// 		MessageBody: messageBody,
// 	}

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		s.logger.Errorf("Failed to marshal SMS payload: %v", err)
// 		return errors.New(localization.ErrorOTPSendFailed.Code)
// 	}

// 	// Create HTTP request
// 	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL, strings.NewReader(string(payloadBytes)))
// 	if err != nil {
// 		s.logger.Errorf("Failed to create HTTP request: %v", err)
// 		return errors.New(localization.ErrorExternalServiceError.Code)
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	// Make request
// 	startTime := time.Now()
// 	resp, err := s.httpClient.Do(req)
// 	duration := time.Since(startTime)

// 	if err != nil {
// 		s.logger.Errorf("SMS API call failed: %v, duration: %v", err, duration)
// 		return errors.New(localization.ErrorExternalServiceError.Code)
// 	}
// 	defer resp.Body.Close()

// 	return nil
// }

// // SendOTP sends an OTP SMS with formatted message
// func (s *SMSPersistence) SendOTP(ctx context.Context, recipient, otpCode string, validMinutes int) error {
// 	messageBody := fmt.Sprintf("Your device lookup OTP is: %s. Valid for %d minutes.", otpCode, validMinutes)
// 	return s.SendSMS(ctx, recipient, messageBody)
// }
