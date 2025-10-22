package budget_category

import (
	"mime/multipart"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var hexColorRegex = regexp.MustCompile(`^#?([a-fA-F\d]{2}){3}$`)

func hasAllowedExtension(filename string, allowed []string) bool {
	if filename == "" {
		return false
	}
	filename = strings.ToLower(strings.TrimSpace(filename))
	for _, ext := range allowed {
		if strings.HasSuffix(filename, strings.ToLower(ext)) {
			return true
		}
	}
	return false
}

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_image_invalid", "invalid image file")
	}
	if file.Size > 10*1024*1024 {
		return validation.NewError("validation_image_size", "image file size must not exceed 10MB")
	}
	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
		return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
	}
	return nil
}

func (b CreateBudgetRequest) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 100).Error("name must be between 1 and 100 characters"),
		),
		validation.Field(&b.Color,
			validation.Required.Error("color is required"),
			validation.Length(7, 7).Error("color must be 7 characters long"),
			validation.Match(hexColorRegex).Error("invalid color format, please enter a valid hex color"),
		),
		validation.Field(&b.Icon,
			validation.Required.Error("icon is required"),
			validation.By(validateImage),
		),
	)
}

func (r UpdateBudgetRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.When(r.Name != nil,
				validation.Length(1, 100).Error("name must be between 1 and 100 characters"),
			),
		),
		validation.Field(&r.Color,
			validation.When(r.Color != nil,
				validation.Length(7, 7).Error("color must be 7 characters long"),
				validation.Match(hexColorRegex).Error("invalid color format, please enter a valid hex color"),
			),
		),
		validation.Field(&r.Icon,
			validation.When(r.Icon != nil,
				validation.By(validateImage),
			),
		),
	)
}
