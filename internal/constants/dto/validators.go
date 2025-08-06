package dto

import (
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func isWeakPin(pin string) bool {
	if strings.Count(pin, string(pin[0])) == len(pin) {
		return true
	}
	isAscending := true
	isDescending := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			isAscending = false
		}
		if pin[i] != pin[i-1]-1 {
			isDescending = false
		}
	}
	return isAscending || isDescending
}

// Validators Group
func (v VerifyForgetPinOtpRequest) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.ResetSessionID, validation.Required),
		validation.Field(&v.Phone, validation.Required, validation.By(func(value interface{}) error {
			phoneStr := fmt.Sprintf("%s", value)
			re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

			matches := re.FindStringSubmatch(phoneStr)
			if matches == nil {
				return fmt.Errorf("invalid phone number")
			}
			return nil
		})),
		validation.Field(&v.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&v.OTP, validation.Required, validation.Length(6, 6)),
	)
}

func (u UpdateProfilePicture) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ProfilePicture,
			validation.Required.Error("profile_picture is required"),
			validation.By(func(value interface{}) error {
				file, ok := value.(*multipart.FileHeader)
				if !ok {
					return fmt.Errorf("invalid file type")
				}

				if file.Size > constants.MaxImageSize {
					return fmt.Errorf("image size exceeds %dMB limit", constants.MaxImageSize)
				}
				if !utils.IsValidImage(file) {
					return fmt.Errorf("invalid image type")
				}
				return nil
			}),
		),
	)
}

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Phone, validation.Required.Error("phone number is required"),
			validation.By(func(value interface{}) error {
				phoneStr := fmt.Sprintf("%s", value)
				re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

				matches := re.FindStringSubmatch(phoneStr)
				if matches == nil {
					return fmt.Errorf("invalid phone number")
				}
				return nil
			})),
		validation.Field(&r.FullName, validation.Required, validation.By(utils.ValidateFullName)),
		validation.Field(&r.Email, validation.When(r.Email != "", is.Email)),
		validation.Field(&r.AppVersion, validation.Required.Error("app version is required")),
	)
}

// Validate validates DeviceLookupRequest
func (d DeviceLookupRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.DeviceUUID,
			validation.Required.Error("device uuid is required"),
			validation.Length(10, 100).Error("device uuid length between 10 and 100")),
		validation.Field(&d.Platform,
			validation.Required.Error("platform is required"),
			validation.In(constants.Android, constants.Ios).Error("platform value should be 'ANDROID', or 'IOS'")),
		validation.Field(&d.AppVersion, validation.Required.Error("app version is required")),
		validation.Field(&d.ApplicationInstallationDate, validation.Required.Error("application installation date is required")),
	)
}

func (c ChangePinRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OldPin, validation.Required, validation.Length(6, 6), is.Digit),
		validation.Field(&c.NewPin, validation.Required, validation.Length(6, 6), is.Digit,
			validation.By(func(value interface{}) error {
				pin := fmt.Sprintf("%s", value)
				if utils.IsWeakPin(pin) {
					return fmt.Errorf("PIN contains weak patterns")
				}
				return nil
			})),
	)
}

func (s SetPinRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.NewPin, validation.Required, validation.Length(6, 6), is.Digit,
			validation.By(func(value interface{}) error {
				pin := fmt.Sprintf("%s", value)
				if utils.IsWeakPin(pin) {
					return errors.New("PIN contains weak patterns")
				}
				return nil
			})),
		validation.Field(&s.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&s.Realm, validation.Required, validation.In(constants.MEMBER_REALM, constants.BANK_REALM, constants.DISTRICT_REALM, constants.BRANCH_REALM, constants.MERCHANT_REALM, constants.COMPANY_REALM)),
	)
}

func (s SetProfileThemeRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ThemeType, validation.Required),
	)
}

func (v VerifyOTPRequest) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.OTP, validation.Required.Error("otp code is required"), validation.Length(6, 6)),
		validation.Field(&v.OtpFor, validation.Required.Error("otp for is required")),
		validation.Field(&v.DeviceUUID, validation.Required.Error("device uuid is required")),
	)
}

func (r ResetPinRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ResetSessionID, validation.Required),
		validation.Field(&r.Phone, validation.Required, validation.By(func(value interface{}) error {
			phoneStr := fmt.Sprintf("%s", value)
			re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

			matches := re.FindStringSubmatch(phoneStr)
			if matches == nil {
				return fmt.Errorf("invalid phone number")
			}
			return nil
		})),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.OTP, validation.Required, validation.Length(6, 6), is.Digit),
		validation.Field(&r.NewPin, validation.Required, validation.Length(6, 6), is.Digit,
			validation.By(func(value interface{}) error {
				pin := fmt.Sprintf("%s", value)
				if utils.IsWeakPin(pin) {
					return errors.New("PIN contains weak patterns")
				}
				return nil
			})),
	)
}

func (p PhoneLoginRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Phone,
			validation.Required.Error("phone is required"),
			validation.By(func(value interface{}) error {
				phoneStr := fmt.Sprintf("%s", value)
				re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

				matches := re.FindStringSubmatch(phoneStr)
				if matches == nil {
					return fmt.Errorf("invalid phone number")
				}
				return nil
			})))
}

func (c CompleteRegistrationRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.RegistrationID, validation.Required),
		validation.Field(&c.Phone, validation.Required, validation.By(func(value interface{}) error {
			phoneStr := fmt.Sprintf("%s", value)
			re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

			matches := re.FindStringSubmatch(phoneStr)
			if matches == nil {
				return fmt.Errorf("invalid phone number")
			}
			return nil
		})),
		validation.Field(&c.FullName, validation.Required),
		validation.Field(&c.OTP, validation.Required),
		validation.Field(&c.Platform, validation.Required, validation.In(constants.Android, constants.Ios)),
		validation.Field(&c.DeviceUUID, validation.Required, validation.Min(10)),
	)

}

func (f ForgetPinSendOtpRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Phone, validation.Required, validation.By(func(value interface{}) error {
			phoneStr := fmt.Sprintf("%s", value)
			re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

			matches := re.FindStringSubmatch(phoneStr)
			if matches == nil {
				return fmt.Errorf("invalid phone number")
			}
			return nil
		})))
}

func (r ResetPinWithTokenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ResetSessionID, validation.Required),
		validation.Field(&r.Phone, validation.Required, validation.By(func(value interface{}) error {
			phoneStr := fmt.Sprintf("%s", value)
			re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

			matches := re.FindStringSubmatch(phoneStr)
			if matches == nil {
				return fmt.Errorf("invalid phone number")
			}
			return nil
		})),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.NewPin, validation.Required, validation.Length(6, 6), is.Digit,
			validation.By(func(value interface{}) error {
				pin := fmt.Sprintf("%s", value)
				if isWeakPin(pin) {
					return fmt.Errorf("PIN contains weak patterns")
				}
				return nil
			})),
	)
}
