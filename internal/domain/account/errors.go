package account

import (
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
     "cbe-super-app-member-users/internal/domain/account/error"
    "cbe-super-app-member-users/internal/shared"
)

type ServiceError struct {
    Code    string
    Message string
}

func (e ServiceError) Error() string {
    return e.Message
}

func NewServiceError(def common.ErrorDefinition) ServiceError {
    return ServiceError{
        Code:    def.Code,
        Message: def.Message,
    }
}

func NewInternalServiceError(def error.ErrorDefinition) ServiceError {
    return ServiceError{
        Code:    def.Code,
        Message: def.Message,
    }
}
func NewlocalServiceError(def shared.ErrorDefinition) ServiceError {
    return ServiceError{
        Code:    def.Code,
        Message: def.Message,
    }
}