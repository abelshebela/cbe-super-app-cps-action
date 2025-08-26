package productcode

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u UpdateProductCodeRequest) Validate() error {
	u.ProductName = strings.TrimSpace(u.ProductName)
	return validation.ValidateStruct(&u,
		validation.Field(&u.ProductName),
		validation.Field(&u.CBEProductCodes,
			validation.By(func(value interface{}) error {
				if pc, ok := value.(model.ProductCodes); ok {
					return pc.Validate()
				}
				return nil
			}),
		),
		validation.Field(&u.CBEIFBProductCodes,
			validation.By(func(value interface{}) error {
				if pc, ok := value.(model.ProductCodes); ok {
					return pc.Validate()
				}
				return nil
			}),
		),
	)
}
