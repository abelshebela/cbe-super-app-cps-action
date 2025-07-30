package unlink

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetUserByAccountRequest struct {
	AccountNumbers string `json:"account_number"`
}

type UnlinkUserRequest struct {
	UserCode string `json:"user_code"`
}

func (un *UnlinkUserRequest) Validate() error {
	err := validation.ValidateStruct(
		un,
		validation.Field(&un.UserCode, validation.Required),
	)
	return err
}
func (ac *GetUserByAccountRequest) Validate() error {
	return validation.ValidateStruct(
		ac,
		validation.Field(&ac.AccountNumbers, validation.Required, validation.Length(10, 20)),
	)
}
