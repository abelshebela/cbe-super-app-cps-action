package encryption

import validation "github.com/go-ozzo/ozzo-validation/v4"

type EncryptionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e *EncryptionRequest) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Username, validation.Required),
		validation.Field(&e.Password, validation.Required),
	)
}
