package roles_dto

import (
	"fmt"
	"strings"

	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var validRoleTypes = []interface{}{"CPS", "BPS"}

func (r CreateJobRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(2, 100).Error("name must be between 2 and 100 characters"),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.Type,
			validation.Required.Error("type is required"),
			validation.In(validRoleTypes...).Error("type must be CPS or BPS"),
		),
		// validation.Field(&r.PortalCards,
		// 	validation.By(validatePortalCards),
		// ),
	)
}

func (r UpdateJobRoleRequest) Validate() error {
	if strings.TrimSpace(r.Code) == "" && strings.TrimSpace(r.Name) == "" && strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("at least one field must be provided for update")
	}
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.When(r.Name != "",
				validation.Length(2, 100).Error("name must be between 2 and 100 characters"),
				validation.By(utils.NoSpecialChars),
			),
		),
		validation.Field(&r.Type,
			validation.When(r.Type != "",
				validation.In(validRoleTypes...).Error("type must be CPS or BPS"),
			),
		),
	)
}

// validatePortalCards validates portal cards array
func validatePortalCards(value interface{}) error {
	cards, ok := value.([]string)
	if !ok {
		return validation.NewError("validation_portal_cards", "portal cards must be an array of strings")
	}

	for i, card := range cards {
		if strings.TrimSpace(card) == "" {
			return validation.NewError("validation_portal_cards", fmt.Sprintf("portal card at index %d cannot be empty", i))
		}

		if err := utils.NoSpecialChars(card); err != nil {
			return validation.NewError("validation_portal_cards", fmt.Sprintf("portal card at index %d contains invalid characters", i))
		}
	}

	return nil
}
