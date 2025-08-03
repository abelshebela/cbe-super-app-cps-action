package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// SMSRequest represents the request for sending SMS
type SMSRequest struct {
	Recipient   string `json:"recipient" validate:"required"`
	MessageBody string `json:"message_body" validate:"required"`
}

// SMSResponse represents the response from SMS API
type SMSResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ExternalCallRequest represents a generic external API call request
type ExternalCallRequest struct {
	URL     string            `json:"url" validate:"required,url"`
	Method  string            `json:"method" validate:"required"`
	Payload string            `json:"payload,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// ExternalCallResponse represents a generic external API call response
type ExternalCallResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body"`
	Error      string            `json:"error,omitempty"`
}

func (r SMSRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Recipient, validation.Required.Error("recipient is required"), is.Digit),
		validation.Field(&r.MessageBody, validation.Required.Error("message body is required")),
	)
}

func (r ExternalCallRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.URL, validation.Required.Error("URL is required"), validation.By(is.URL)),
		validation.Field(&r.Method, validation.Required.Error("method is required"), validation.In("GET", "POST", "PUT", "DELETE", "PATCH")),
	)
}
