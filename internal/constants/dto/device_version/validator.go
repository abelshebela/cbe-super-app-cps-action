package deviceversion

import (
	"regexp"
	"strings"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Clean normalizes fields for create
func (r *CreateDeviceVersionRequest) Clean() {
	r.Platform = strings.TrimSpace(strings.ToLower(r.Platform))
	r.LatestVersion = strings.TrimSpace(r.LatestVersion)
	r.ReleaseNotes = strings.TrimSpace(r.ReleaseNotes)
}

func (r CreateDeviceVersionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.LatestVersion,
			validation.Required,
			validation.Length(1, 50),
			validation.Match(regexp.MustCompile(`^[0-9.]+$`)).Error("LatestVersion must contain only digits and dots (0-9, .)"),
		),
		validation.Field(&r.Platform,
			validation.Required,
			validation.In("ANDROID", "IOS", "android", "ios"),
		),
		validation.Field(&r.ReleaseNotes,
			validation.Length(0, 500),
		),
	)
}

// Clean normalizes fields for update
func (r *UpdateDeviceVersionRequest) Clean() {
	r.Platform = strings.TrimSpace(strings.ToLower(r.Platform))
	r.LatestVersion = strings.TrimSpace(r.LatestVersion)
}

func (r UpdateDeviceVersionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID,
			validation.Required,
			is.Hexadecimal,
			validation.Length(24, 24),
		),
		validation.Field(&r.LatestVersion,
			validation.Length(1, 50),
		),
		validation.Field(&r.Platform,
			validation.In("ANDROID", "IOS", "android", "ios"),
		),
		validation.Field(&r.ReleaseNotes,
			validation.Length(0, 500),
		),
	)
}

func (r *EnableOrDisableDeviceVersion) Clean() {
	r.ID = strings.TrimSpace(r.ID)
}

func (r EnableOrDisableDeviceVersion) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID,
			validation.Required,
			is.Hexadecimal,
			validation.Length(24, 24),
		),
	)
}
