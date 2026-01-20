package bpsuser

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// type FullName struct {
// 	FirstName  string `json:"first_name"`
// 	MiddleName string `json:"middle_name"`
// 	LastName   string `json:"last_name"`
// }

type BPSUserCreateRequest struct {
	UserID      string `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	JobTitle    string `json:"job_title"`
	// Role        string   `json:"role"`
	BranchCode []string `json:"branch_code"`
	// HomeBranch  string   `json:"home_branch"`
	// Realm       string   `json:"realm"`
	// Enabled bool `json:"enabled"`
}

type BPSUserUpdateRequest struct {
	UserID      string `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	JobTitle    string `json:"job_title"`
	// Role        string   `json:"role"`
	BranchCode []string `json:"branch_code"`
	// HomeBranch  string   `json:"home_branch"`
	// Realm       string   `json:"realm"`
	// Enabled bool `json:"enabled"`
}

var phoneRegex = regexp.MustCompile(`^\+251[97]\d{8}$`)

func (r *BPSUserCreateRequest) Validate() error {
	// Trim spaces
	r.UserID = strings.TrimSpace(r.UserID)
	r.FullName = strings.TrimSpace(r.FullName)
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
	r.Email = strings.TrimSpace(r.Email)
	r.JobTitle = strings.TrimSpace(r.JobTitle)

	return validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Match(phoneRegex).Error("must be a valid Ethiopian phone number (+2519xxxxxxxx or +2517xxxxxxxx)")),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.JobTitle, validation.Required),
		validation.Field(&r.BranchCode, validation.Required),
	)
}

func (r *BPSUserUpdateRequest) Validate() error {
	// Trim spaces
	r.UserID = strings.TrimSpace(r.UserID)
	r.FullName = strings.TrimSpace(r.FullName)
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
	r.Email = strings.TrimSpace(r.Email)
	r.JobTitle = strings.TrimSpace(r.JobTitle)

	return validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Match(phoneRegex).Error("must be a valid Ethiopian phone number (+2519xxxxxxxx or +2517xxxxxxxx)")),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.JobTitle, validation.Required),
		validation.Field(&r.BranchCode, validation.Required),
	)
}
