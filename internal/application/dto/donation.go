package dto

import (
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type DonationCategoryRequest struct {
	CategoryName string                `json:"category_name"`
	Icon         *multipart.FileHeader `json:"donation_icon,omitempty"`
}

type DonationCategoryResponse struct {
	CategoryName string `json:"category_name"`
	Icon         string `json:"donation_icon"`
}

type DonationCategoryListResponse struct {
	CategoryName   string `json:"category_name"`
	Icon           string `json:"donation_icon"`
	IsDeleted      bool   `json:"is_deleted"`
	CreatedAt      string `json:"created_at"`
	LastModifiedAt string `json:"last_modified_at"`
}

// DonationCategoryCPSRequest is used for CPS actions where icon is a URL string
type DonationCategoryCPSRequest struct {
	CategoryName string `json:"category_name"`
	Icon         string `json:"donation_icon,omitempty"`
}

func (d DonationCategoryRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CategoryName,
			validation.Required.Error("category name is required"),
			validation.Length(3, 100).Error("category name must be between 3 and 100 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-&]+$`)).Error("category name can only contain letters, numbers, spaces, hyphens, and ampersands"),
		),
		validation.Field(&d.Icon,
			validation.Required.Error("category icon is required"),
			validation.By(validateDonationIcon),
		),
	)
}

func (d DonationCategoryRequest) ValidateForUpdate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CategoryName,
			validation.Required.Error("category name is required"),
			validation.Length(3, 100).Error("category name must be between 3 and 100 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-&]+$`)).Error("category name can only contain letters, numbers, spaces, hyphens, and ampersands"),
		),
		validation.Field(&d.Icon,
			validation.When(d.Icon != nil, validation.By(validateDonationIcon)),
		),
	)
}

type DonationCompanyRequest struct {
	CompanyName   string                `json:"company_name"`
	CompanyLogo   *multipart.FileHeader `json:"company_logo"`
	AccountNumber string                `json:"account_number"`
}

func (d DonationCompanyRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyName,
			validation.Required.Error("company name is required"),
			validation.Length(2, 200).Error("company name must be between 2 and 200 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-&.]+$`)).Error("company name can only contain letters, numbers, spaces, hyphens, periods, and ampersands"),
		),
		validation.Field(&d.CompanyLogo,
			validation.Required.Error("company logo is required"),
			validation.By(validateLogo),
		),
		validation.Field(&d.AccountNumber,
			validation.Required.Error("account number is required"),
			validation.Length(8, 50).Error("account number must be between 8 and 50 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)).Error("account number can only contain letters, numbers, and hyphens"),
		),
	)
}

type DonationRequest struct {
	CompanyID           string                  `json:"company_id"`
	CategoryID          string                  `json:"category_id"`
	Title               string                  `json:"title"`
	IsFeatured          bool                    `json:"is_featured"`
	Target              float64                 `json:"donation_amount"`
	DonationDescription string                  `json:"donation_description"`
	DonationImages      []*multipart.FileHeader `json:"donation_images"`
	EndDate             time.Time               `json:"end_date"`
	StartDate           time.Time               `json:"start_date"`
}

func (d DonationRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyID,
			validation.Required.Error("company ID is required"),
			is.UUID.Error("company ID must be a valid UUID"),
		),
		validation.Field(&d.CategoryID,
			validation.Required.Error("category ID is required"),
			is.UUID.Error("category ID must be a valid UUID"),
		),
		validation.Field(&d.Title,
			validation.Required.Error("title is required"),
			validation.Length(5, 200).Error("title must be between 5 and 200 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-!?.&]+$`)).Error("title can only contain letters, numbers, spaces, and basic punctuation"),
		),
		validation.Field(&d.Target,
			validation.Required.Error("donation amount is required"),
			validation.Min(1.0).Error("donation amount must be at least 1"),
			validation.Max(100000000.0).Error("donation amount must not exceed 100,000,000"),
		),
		validation.Field(&d.DonationDescription,
			validation.Required.Error("donation description is required"),
			validation.Length(10, 5000).Error("description must be between 10 and 5000 characters"),
		),
		validation.Field(&d.DonationImages,
			validation.Required.Error("at least one donation image is required"),
			validation.Length(1, 10).Error("must provide between 1 and 10 images"),
			validation.Each(validation.By(validateDonationImage)),
		),
		validation.Field(&d.StartDate,
			validation.Required.Error("start date is required"),
			validation.By(validateStartDate),
		),
		validation.Field(&d.EndDate,
			validation.Required.Error("end date is required"),
			validation.By(validateEndDate(d.StartDate)),
		),
	)
}
func validateDonationIcon(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_icon_invalid", "invalid icon file")
	}
	if file.Size > 5*1024*1024 {
		return validation.NewError("validation_icon_size", "icon file size must not exceed 5MB")
	}
	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
		return validation.NewError("validation_icon_format", "icon must be JPG, JPEG, PNG, or GIF")
	}
	return nil
}

func validateLogo(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_logo_invalid", "invalid logo file")
	}
	if file.Size > 5*1024*1024 {
		return validation.NewError("validation_logo_size", "logo file size must not exceed 5MB")
	}
	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
		return validation.NewError("validation_logo_format", "logo must be JPG, JPEG, PNG, or GIF")
	}
	return nil
}

func validateDonationImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_image_invalid", "invalid image file")
	}
	if file.Size > 10*1024*1024 {
		return validation.NewError("validation_image_size", "each image file size must not exceed 10MB")
	}
	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png"}) {
		return validation.NewError("validation_image_format", "each image must be JPG, JPEG, or PNG")
	}
	return nil
}

func validateStartDate(value interface{}) error {
	startDate, ok := value.(time.Time)
	if !ok {
		return validation.NewError("validation_start_date_invalid", "invalid start date")
	}
	if startDate.Before(time.Now().AddDate(0, 0, -1)) {
		return validation.NewError("validation_start_date_past", "start date cannot be in the past")
	}
	return nil
}

func validateEndDate(startDate time.Time) validation.RuleFunc {
	return func(value interface{}) error {
		endDate, ok := value.(time.Time)
		if !ok {
			return validation.NewError("validation_end_date_invalid", "invalid end date")
		}
		if !startDate.IsZero() && endDate.Before(startDate) {
			return validation.NewError("validation_end_date_before_start", "end date must be after start date")
		}
		if endDate.Before(time.Now()) {
			return validation.NewError("validation_end_date_past", "end date cannot be in the past")
		}
		return nil
	}
}

func hasAllowedExtension(filename string, allowed []string) bool {
	filename = strings.ToLower(filename)
	for _, ext := range allowed {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
