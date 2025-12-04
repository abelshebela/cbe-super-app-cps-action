package department_dto

import (
	"fmt"
	"strings"

	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r CreateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Department,
			validation.Required.Error("department name is required"),
			validation.Length(2, 100).Error("department name must be between 2 and 100 characters"),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.PortalCards,
			validation.By(validatePortalCards),
		),
	)
}

func (r UpdateDepartmentRequest) Validate() error {
	if r.Department == "" && len(r.PortalCards) == 0 {
		return fmt.Errorf("at least one field must be provided for update")
	}

	return validation.ValidateStruct(&r,
		validation.Field(&r.Department,
			validation.When(r.Department != "",
				validation.Length(2, 100).Error("department name must be between 2 and 100 characters"),
				validation.By(utils.NoSpecialChars),
			),
		),
		validation.Field(&r.PortalCards,
			validation.When(len(r.PortalCards) > 0,
				validation.By(validatePortalCards),
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
