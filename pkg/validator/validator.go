package validator

import (
	"cbe-super-app-member-auth/internal/adapter/inbound/dto"
	"cbe-super-app-member-auth/pkg/utils"
	"regexp"
	"strconv"
	"strings"
)

func LoginInputValidator(input dto.UsernameLoginRequest) []string {
	var errs []string

	if strings.TrimSpace(input.UserName) == "" {
		errs = append(errs, "Username is required")
	}
	err := utils.ValidateInputNoSpecialChars(input.UserName)

	if errs == nil {
		if err != nil {
			errs = append(errs, "Invalid Input")
		}
	}
	return errs
}

func ResetPassInputValidator(input dto.PassResetRequest) []string {
	var errs []string

	if strings.TrimSpace(input.UserName) == "" {
		errs = append(errs, "Username is required")
		err := utils.ValidateInputNoSpecialChars(input.UserName)

		if errs == nil {
			if err != nil {
				errs = append(errs, "Invalid Input")
			}
		}
	}

	if strings.TrimSpace(input.NewPassword) == "" {
		errs = append(errs, "Password is required")

	} else if len(input.NewPassword) < 6 {
		errs = append(errs, "Password must be at least 6 characters")
	}

	return errs
}

func ValidatePhoneNumber(phone string) bool {
	var re = regexp.MustCompile(`^(?:\+251|251|0)?9\d{8}$`)
	return re.MatchString(phone)
}
func PhoneLookUpValidator(input dto.PhoneLookUpRequest) []string {
	var errs []string

	if strings.TrimSpace(input.PhoneNumber) == "" {
		errs = append(errs, "Phone number is required")
	} else {
		phone := ValidatePhoneNumber(input.PhoneNumber)

		if !phone {
			errs = append(errs, "Invalid Phone Input")
		}
	}

	return errs
}

func MemberSelfSignupValidator(input dto.MemberSelfSignup) []string {
	var errs []string

	if strings.TrimSpace(input.FullName) == "" {
		errs = append(errs, "Full name is required")
	} else {
		err := utils.ValidateInputNoSpecialChars(input.FullName)

		if errs == nil {
			if err != nil {
				errs = append(errs, "Invalid Input on Full name")
			}
		}
	}

	if strings.TrimSpace(input.PhoneNumber) == "" {
		errs = append(errs, "Phone number is required")
	} else {
		phone := ValidatePhoneNumber(input.PhoneNumber)

		if !phone {
			errs = append(errs, "Invalid Phone Input")
		}
	}

	if strings.TrimSpace(string(input.Gender)) == "" {
		errs = append(errs, "Gender is required")
	} else {
		err := utils.ValidateInputNoSpecialChars(string(input.Gender))

		if errs == nil {
			if err != nil {
				errs = append(errs, "Invalid Input Gender")
			}
		}
	}

	if strings.TrimSpace(input.BirthDate) == "" {
		errs = append(errs, "Birth date is required")
	}

	if strings.TrimSpace(input.Email) == "" {
		errs = append(errs, "Email is required")
	}

	if strings.TrimSpace(input.City) == "" {
		errs = append(errs, "City is required")
	} else {
		err := utils.ValidateInputNoSpecialChars(string(input.City))

		if errs == nil {
			if err != nil {
				errs = append(errs, "Invalid Input city")
			}
		}
	}

	return errs
}

func PinStrengthValidator(input dto.PinStrengthRequest) []string {
	var errs []string

	if strings.TrimSpace(input.PhoneNumber) == "" {
		errs = append(errs, "Phone number is required")
	} else {
		phone := ValidatePhoneNumber(input.PhoneNumber)

		if !phone {
			errs = append(errs, "Invalid Phone Input")
		}
	}

	pin := strconv.Itoa(input.NewPin)

	if strings.TrimSpace(pin) == "0" {
		errs = append(errs, "Pin is required")
	} else if len(pin) < 6 || len(pin) > 6 {
		errs = append(errs, "Pin 2 must be only 6 characters")
	}

	return errs
}

func PinCodeValidator(input dto.SetLoginPINRequest) []string {
	var errs []string

	if strings.TrimSpace(input.PhoneNumber) == "" {
		errs = append(errs, "Phone number is required")
	} else {
		phone := ValidatePhoneNumber(input.PhoneNumber)

		if !phone {
			errs = append(errs, "Invalid Phone Input")
		}
	}

	if strings.TrimSpace(input.PinCode) == "" {
		errs = append(errs, "Pin is required")
	} else if len(input.PinCode) < 6 || len(input.PinCode) > 6 {
		errs = append(errs, "Password must be only 6 characters")
	}

	return errs
}
func LoginPassInputValidator(input dto.LoginWithPassRequest) []string {
	var errs []string

	if strings.TrimSpace(input.UserName) == "" {
		errs = append(errs, "Username is required")
	} else {
		err := utils.ValidateInputNoSpecialChars(input.UserName)

		if errs == nil {
			if err != nil {
				errs = append(errs, "Invalid Input")
			}
		}
	}

	if strings.TrimSpace(input.Password) == "" {
		errs = append(errs, "Password is required")
	} else if len(input.Password) < 6 {
		errs = append(errs, "Password must be at least 6 characters")
	}

	return errs
}
func ChangePassInputValidator(input dto.ChangePasswordRequest) []string {
	var errs []string

	if strings.TrimSpace(input.Username) == "" {
		errs = append(errs, "Username is required")
	} else {
		err := utils.ValidateInputNoSpecialChars(input.Username)

		if errs == nil {
			if err != nil {
				errs = append(errs, "Invalid Input")
			}
		}
	}

	if strings.TrimSpace(input.OldPassword) == "" {
		errs = append(errs, "Password is required")
	} else if len(input.NewPassword) < 6 {
		errs = append(errs, "Password must be at least 6 characters")
	} else if input.NewPassword == input.OldPassword {
		errs = append(errs, "You entered the same password")
	}

	return errs
}

func OtpApproveInputValidator(input dto.OtpApproveRequest) []string {
	var errs []string

	if strings.TrimSpace(input.UserName) == "" {
		errs = append(errs, "Username is required")
	}

	if strings.TrimSpace(input.OtpCode) == "" {
		errs = append(errs, "Password is required")
	}

	if len(strings.TrimSpace(input.OtpCode)) < 6 || len(strings.TrimSpace(input.OtpCode)) > 6 {
		errs = append(errs, "The Otp code not have to be greater or less than 6")
	}

	_, err := strconv.Atoi(input.OtpCode)
	if err != nil {
		errs = append(errs, "Otp Code must be only number")
	}

	return errs
}
