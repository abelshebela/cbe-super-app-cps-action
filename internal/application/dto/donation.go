package dto

import (
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type DonationImage struct {
	ID        string `json:"id" bson:"id"`
	PhotoURL  string `json:"photo_url" bson:"photo_url"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

type DonationCategoryRequest struct {
	CategoryName string                `json:"category_name" bson:"category_name"`
	Icon         *multipart.FileHeader `json:"donation_icon,omitempty" bson:"donation_icon,omitempty"`
}

type DonationCategoryResponse struct {
	CategoryName string `json:"category_name" bson:"category_name"`
	Icon         string `json:"donation_icon" bson:"donation_icon"`
}

// DonationCategoryListResponse is used for fetching donation categories
type DonationCategoryListResponse struct {
	ID             string `json:"id" bson:"id"`
	CategoryName   string `json:"category_name" bson:"category_name"`
	Icon           string `json:"donation_icon" bson:"donation_icon"`
	IsDeleted      bool   `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
}

// DonationCategoryCPSRequest is used for CPS actions where icon is a URL string
type DonationCategoryCPSRequest struct {
	ID           string `json:"id,omitempty" bson:"id,omitempty"`
	CategoryName string `json:"category_name" bson:"category_name"`
	Icon         string `json:"donation_icon,omitempty" bson:"donation_icon,omitempty"`
}

// DonationCompanyResponse is used for service layer responses
type DonationCompanyResponse struct {
	CompanyName   string `json:"company_name" bson:"company_name"`
	CompanyLogo   string `json:"company_logo" bson:"company_logo"`
	AccountNumber string `json:"account_number" bson:"account_number"`
}

// DonationCompanyListResponse is used for fetching donation companies
type DonationCompanyListResponse struct {
	ID             string `json:"id" bson:"id"`
	CompanyName    string `json:"company_name" bson:"company_name"`
	CompanyLogo    string `json:"company_logo" bson:"company_logo"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	IsDeleted      bool   `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
}

// DonationCompanyCPSRequest is used for CPS actions where logo is a URL string
type DonationCompanyCPSRequest struct {
	ID            string `json:"id,omitempty" bson:"id,omitempty"`
	CompanyName   string `json:"company_name" bson:"company_name"`
	CompanyLogo   string `json:"company_logo,omitempty" bson:"company_logo,omitempty"`
	AccountNumber string `json:"account_number" bson:"account_number"`
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
			validation.By(validateImage),
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
			validation.When(d.Icon != nil, validation.By(validateImage)),
		),
	)
}

type DonationCompanyRequest struct {
	CompanyName   string                `json:"company_name" bson:"company_name"`
	CompanyLogo   *multipart.FileHeader `json:"company_logo" bson:"company_logo"`
	AccountNumber string                `json:"account_number" bson:"account_number"`
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
			validation.By(validateImage),
		),
		validation.Field(&d.AccountNumber,
			validation.Required.Error("account number is required"),
			validation.Length(8, 50).Error("account number must be between 8 and 50 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)).Error("account number can only contain letters, numbers, and hyphens"),
		),
	)
}

func (d DonationCompanyRequest) ValidateForUpdate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyName,
			validation.Required.Error("company name is required"),
			validation.Length(2, 200).Error("company name must be between 2 and 200 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-&.]+$`)).Error("company name can only contain letters, numbers, spaces, hyphens, periods, and ampersands"),
		),
		validation.Field(&d.CompanyLogo,
			validation.When(d.CompanyLogo != nil, validation.By(validateImage)),
		),
		validation.Field(&d.AccountNumber,
			validation.Required.Error("account number is required"),
			validation.Length(8, 50).Error("account number must be between 8 and 50 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)).Error("account number can only contain letters, numbers, and hyphens"),
		),
	)
}

// DonationResponse is used for service layer responses
type DonationResponse struct {
	DonationCode        string          `json:"donation_code" bson:"donation_code"`
	CompanyID           string          `json:"company_id" bson:"company_id"`
	CategoryID          string          `json:"category_id" bson:"category_id"`
	Title               string          `json:"title" bson:"title"`
	IsFeatured          bool            `json:"is_featured" bson:"is_featured"`
	Target              int             `json:"target" bson:"target"`
	DonationDescription string          `json:"donation_description" bson:"donation_description"`
	DonationImages      []DonationImage `json:"donation_images" bson:"donation_images"`
	CoverImage          string          `json:"cover_image,omitempty" bson:"cover_image,omitempty"` // URL for cover image
	EndDate             string          `json:"end_date" bson:"end_date"`
	StartDate           string          `json:"start_date" bson:"start_date"`
}

// DonationListResponse is used for fetching donations
type DonationListResponse struct {
	ID                  string `json:"id" bson:"id"`
	DonationCode        string `json:"donation_code" bson:"donation_code"`
	Company             DonationCompanyCPSRequest
	Category            DonationCategoryCPSRequest
	Title               string          `json:"title" bson:"title"`
	IsFeatured          bool            `json:"is_featured" bson:"is_featured"`
	Target              int             `json:"target" bson:"target"`
	DonationDescription string          `json:"donation_description" bson:"donation_description"`
	DonationImages      []DonationImage `json:"donation_images" bson:"donation_images"`
	CoverImage          string          `json:"cover_image,omitempty" bson:"cover_image,omitempty"` // URL for cover image
	EndDate             string          `json:"end_date" bson:"end_date"`
	StartDate           string          `json:"start_date" bson:"start_date"`
	IsDeleted           bool            `json:"is_deleted" bson:"is_deleted"`
	CreatedAt           string          `json:"created_at" bson:"created_at"`
	LastModifiedAt      string          `json:"last_modified_at" bson:"last_modified_at"`
}

// DonationCPSRequest is used for CPS actions where images are URL strings
type DonationCPSRequest struct {
	DonationCode        string          `json:"donation_code,omitempty" bson:"donation_code,omitempty"`
	CompanyID           string          `json:"company_id" bson:"company_id"`
	CategoryID          string          `json:"category_id" bson:"category_id"`
	Title               string          `json:"title" bson:"title"`
	IsFeatured          bool            `json:"is_featured" bson:"is_featured"`
	Target              int             `json:"target" bson:"target"`
	DonationDescription string          `json:"donation_description" bson:"donation_description"`
	DonationImages      []DonationImage `json:"donation_images,omitempty" bson:"donation_images,omitempty"`
	CoverImage          string          `json:"cover_image,omitempty" bson:"cover_image,omitempty"`
	EndDate             string          `json:"end_date" bson:"end_date"`
	StartDate           string          `json:"start_date" bson:"start_date"`
	ImageIDToDelete     string          `json:"image_id_to_delete,omitempty" bson:"image_id_to_delete,omitempty"`
}

type DonationRequest struct {
	DonationCode        string                  `json:"donation_code,omitempty" bson:"donation_code,omitempty"`
	CompanyID           string                  `json:"company_id" bson:"company_id"`
	CategoryID          string                  `json:"category_id" bson:"category_id"`
	Title               string                  `json:"title" bson:"title"`
	IsFeatured          bool                    `json:"is_featured" bson:"is_featured"`
	Target              int                     `json:"target" bson:"target"`
	DonationDescription string                  `json:"donation_description" bson:"donation_description"`
	DonationImages      []*multipart.FileHeader `json:"donation_images" bson:"donation_images"`
	CoverImage          *multipart.FileHeader   `json:"cover_image,omitempty" bson:"cover_image,omitempty"`
	EndDate             time.Time               `json:"end_date" bson:"end_date"`
	StartDate           time.Time               `json:"start_date" bson:"start_date"`
}

func (d DonationRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyID,
			validation.Required.Error("company ID is required"),
		),
		validation.Field(&d.CategoryID,
			validation.Required.Error("category ID is required"),
		),
		validation.Field(&d.Title,
			validation.Required.Error("title is required"),
			validation.Length(5, 200).Error("title must be between 5 and 200 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-!?.&]+$`)).Error("title can only contain letters, numbers, spaces, and basic punctuation"),
		),
		validation.Field(&d.Target,
			validation.Required.Error("donation amount is required"),
			validation.By(validateDonationAmount),
		),
		validation.Field(&d.DonationDescription,
			validation.Required.Error("donation description is required"),
		),
		validation.Field(&d.DonationImages,
			validation.When(d.DonationImages != nil, validation.By(validateDonationImages)),
		),
		validation.Field(&d.CoverImage,
			validation.When(d.CoverImage != nil, validation.By(validateImage)),
		),

		validation.Field(&d.StartDate,
			validation.When(!d.StartDate.IsZero(), validation.By(validateStartDate)),
		),
		validation.Field(&d.EndDate,
			validation.Required.Error("end date is required"),
			validation.By(validateEndDate(d.StartDate)),
		),
	)
}

func (d DonationRequest) ValidateForUpdate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyID,
			validation.When(d.CompanyID != "", validation.Required.Error("company ID is required"),
				is.UUID.Error("company ID must be a valid UUID")),
		),
		validation.Field(&d.CategoryID,
			validation.When(d.CategoryID != "", validation.Required.Error("category ID is required"),
				is.UUID.Error("category ID must be a valid UUID")),
		),
		validation.Field(&d.Title,
			validation.When(d.Title != "", validation.Required.Error("title is required"),
				validation.Length(5, 200).Error("title must be between 5 and 200 characters"),
				validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s\-!?.&]+$`)).Error("title can only contain letters, numbers, spaces, and basic punctuation")),
		),
		validation.Field(&d.Target,
			validation.When(d.Target > 0, validation.By(validateDonationAmount)),
		),
		validation.Field(&d.DonationDescription,
			validation.When(d.DonationDescription != "", validation.Required.Error("donation description is required")),
		),

		validation.Field(&d.StartDate,
			validation.When(!d.StartDate.IsZero(), validation.By(validateStartDate)),
		),
		validation.Field(&d.EndDate,
			validation.When(!d.EndDate.IsZero(), validation.Required.Error("end date is required"),
				validation.By(validateEndDate(d.StartDate))),
		),
	)
}

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_image_invalid", "invalid image file")
	}

	// Check file size
	if file.Size > 10*1024*1024 {
		return validation.NewError("validation_image_size", "image file size must not exceed 10MB")
	}

	// Check file extension
	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
		return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
	}

	return nil
}

func validateDonationImages(value interface{}) error {
	files, ok := value.([]*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_images_invalid", "invalid image files")
	}

	for _, file := range files {
		// Check file size
		if file.Size > 10*1024*1024 {
			return validation.NewError("validation_image_size",
				"image exceeds the 10MB size limit")
		}

		// Check file extension
		if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
			return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
		}
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

func validateDonationAmount(value interface{}) error {
	amount, ok := value.(int)
	if !ok {
		return validation.NewError("validation_amount_invalid", "invalid donation amount")
	}

	if amount <= 0 {
		return validation.NewError("validation_amount_zero", "donation amount must be greater than zero")
	}

	if amount > 100000000 {
		return validation.NewError("validation_amount_too_large", "donation amount must not exceed 100,000,000")
	}

	return nil
}

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
