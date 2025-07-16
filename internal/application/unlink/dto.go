package unlink

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var noSpecialChars = validation.Match(regexp.MustCompile(`^[A-Za-z0-9]+$`)).Error("userCode must not contain special characters")

type UnlinkDeviceRequest struct {
	UserCode string `json:"userCode"`
}

func (r UnlinkDeviceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserCode,
			validation.Required,
			noSpecialChars,
		),
	)
}

type ApproveUnlinkDeviceRequest struct {
	UserCode string `json:"userCode"`
	Decision string `json:"decision"`
	Reason   string `json:"rejected_reason"`
}

func (r ApproveUnlinkDeviceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserCode,
			validation.Required,
			noSpecialChars,
		),
		validation.Field(&r.Decision,
			validation.Required.Error("decision is required"),
			validation.In("AUTHORIZED", "DENIED").Error("decision must be 'AUTHORIZED' or 'DENIED'"),
		),
		validation.Field(&r.Reason,
			validation.When(r.Decision == "DENIED", validation.Required.Error("reason is required when declining")),
		),
	)
}

type Response struct {
	Message string `json:"message"`
}
