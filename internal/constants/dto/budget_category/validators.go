package budget_category

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var hexColorRegex = regexp.MustCompile(`^#?([a-fA-F\d]{2}){3}$`)

var validBudgetCategoryTypes = []interface{}{
	string(local_model.BudgetCategoryTypeCB),
	string(local_model.BudgetCategoryTypeIFB),
	string(local_model.BudgetCategoryTypeBOTH),
}

// func hasAllowedExtension(filename string, allowed []string) bool {
// 	if filename == "" {
// 		return false
// 	}
// 	filename = strings.ToLower(strings.TrimSpace(filename))
// 	for _, ext := range allowed {
// 		if strings.HasSuffix(filename, strings.ToLower(ext)) {
// 			return true
// 		}
// 	}
// 	return false
// }

func validateBudgetIcon(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidImage
	}
	if !utils.IsValidImage(file) {
		return localization.ErrorMissingOrInvalidImage
	}

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}

func (b CreateBudgetRequest) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 50).Error("name must be between 1 and 50 characters"),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&b.Color,
			validation.Required.Error("color is required"),
			validation.Length(7, 7).Error("color must be 7 characters long"),
			validation.Match(hexColorRegex).Error("invalid color format, please enter a valid hex color"),
		),
		validation.Field(&b.Icon, validation.By(func(value interface{}) error { return validateBudgetIcon(value) })),
		validation.Field(&b.Type,
			validation.Required.Error("type is required"),
			validation.In(validBudgetCategoryTypes...).Error("type must be one of: CB, IFB, BOTH"),
		),
	)
}

func (r UpdateBudgetRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.When(r.Name != "",
				validation.Length(1, 50).Error("name must be between 1 and 50 characters"),
				validation.By(utils.NoSpecialChars),
			),
		),
		validation.Field(&r.Color,
			validation.When(r.Color != "",
				validation.Length(7, 7).Error("color must be 7 characters long"),
				validation.Match(hexColorRegex).Error("invalid color format, please enter a valid hex color"),
			),
		),
		validation.Field(&r.Icon,
			validation.When(r.Icon != nil,
				validation.By(func(value interface{}) error { return validateBudgetIcon(value) }),
			),
		),
		validation.Field(&r.Type,
			validation.When(r.Type != "",
				validation.In(validBudgetCategoryTypes...).Error("type must be one of: CB, IFB, BOTH"),
			),
		),
	)
}
