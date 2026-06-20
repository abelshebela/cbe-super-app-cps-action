package deviceversion

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"html"
	"regexp"
	"strings"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const errDeviceStateInvalid = "device_state must be one of FORCE_UPDATE, MAINTENANCE, STABLE"

// Clean normalizes fields for create
func (r *CreateDeviceVersionRequest) Clean() {
	r.Platform = strings.TrimSpace(strings.ToLower(r.Platform))
	r.LatestVersion = strings.TrimSpace(r.LatestVersion)
	r.ReleaseNotes = strings.TrimSpace(r.ReleaseNotes)
	r.Platform = html.EscapeString(r.Platform)
	r.ReleaseNotes = html.EscapeString(r.ReleaseNotes)
	r.LatestVersion = html.EscapeString(r.LatestVersion)
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
		validation.Field(&r.DeviceState,
			validation.Required,
			validation.In("FORCE_UPDATE", "MAINTENANCE", "STABLE").Error(errDeviceStateInvalid),
		),
		validation.Field(&r.ReleaseNotes,
			validation.Length(0, 500),
			validation.By(utils.NoSpecialChars),
			validation.Match(regexp.MustCompile(`[0-9a-zA-Z\s.,<>!?'"()-]*`)).Error("ReleaseNotes contains invalid characters"),
		),
	)
}

// Clean normalizes fields for update
func (r *UpdateDeviceVersionRequest) Clean() {
	r.Platform = strings.TrimSpace(strings.ToLower(r.Platform))
	r.LatestVersion = strings.TrimSpace(r.LatestVersion)
	r.ReleaseNotes = strings.TrimSpace(r.ReleaseNotes)
}

func (r UpdateDeviceVersionRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ID,
			validation.Required,
			is.Hexadecimal,
			validation.By(utils.NoSpecialChars),
			validation.Length(24, 24),
		),
		validation.Field(&r.LatestVersion,
			validation.Length(1, 50),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.Platform,
			validation.In("ANDROID", "IOS", "android", "ios").Error("platform must be either of ANDROID or IOS"),
		),
		validation.Field(&r.DeviceState,
			validation.In("FORCE_UPDATE", "MAINTENANCE", "STABLE", "").Error(errDeviceStateInvalid),
		),
		validation.Field(&r.ReleaseNotes,
			validation.Length(0, 500),
			validation.By(utils.NoSpecialChars),
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r SetDeviceStateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DeviceState,
			validation.Required,
			validation.In("FORCE_UPDATE", "MAINTENANCE", "STABLE").Error(errDeviceStateInvalid),
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
