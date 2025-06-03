package users

import (
	"fmt"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
)

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewServiceError(def common.ErrorDefinition) error {
	return &ServiceError{
		Code:    def.Code,
		Message: def.Message,
	}
}