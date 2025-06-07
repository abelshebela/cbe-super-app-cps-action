package users

import (
	"fmt"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
)

type ServiceError struct {
	Code    string
	Message string
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewServiceError(def common.ErrorDefinition) ServiceError {
	return ServiceError{
		Code:    def.Code,
		Message: def.Message,
	}
}

func NewInvalidOTPError() ServiceError {
	return ServiceError{
		Code:    "INVALID_OTP",
		Message: "The provided OTP is invalid",
	}
}

func NewExpiredOTPError() ServiceError {
	return ServiceError{
		Code:    "EXPIRED_OTP",
		Message: "The OTP has expired",
	}
}