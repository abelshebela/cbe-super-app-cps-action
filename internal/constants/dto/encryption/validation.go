package encryption

import (
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (e *EncryptionRequest) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Username, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&e.Password, validation.Required),
	)
}
