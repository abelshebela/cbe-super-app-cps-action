package bpsuser

import (
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	UserID      *string `json:"user_id"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	JobTitle    *string `json:"job_title"`
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
		validation.Field(&r.UserID, validation.Required, validation.Length(3, 10)),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Match(phoneRegex).Error("must be a valid Ethiopian phone number (+2519xxxxxxxx or +2517xxxxxxxx)")),
		validation.Field(&r.Email, validation.Required, is.Email, validation.Match(regexp.MustCompile(`^[A-Za-z0-9._-]+@cbe\.com\.et$`)).Error("must be a valid @cbe.com.et email")),
		validation.Field(&r.JobTitle, validation.Required),
		validation.Field(&r.BranchCode, validation.Required),
	)
}

func (r *BPSUserUpdateRequest) Validate() error {
	// Trim spaces for non-nil pointers
	if r.UserID != nil {
		*r.UserID = strings.TrimSpace(*r.UserID)
	}
	if r.FullName != nil {
		*r.FullName = strings.TrimSpace(*r.FullName)
	}
	if r.PhoneNumber != nil {
		*r.PhoneNumber = strings.TrimSpace(*r.PhoneNumber)
	}
	if r.Email != nil {
		*r.Email = strings.TrimSpace(*r.Email)
	}
	if r.JobTitle != nil {
		*r.JobTitle = strings.TrimSpace(*r.JobTitle)
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required, validation.Length(3, 10)),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Match(phoneRegex).Error("must be a valid Ethiopian phone number (+2519xxxxxxxx or +2517xxxxxxxx)")),
		validation.Field(&r.Email, validation.Required, is.Email, validation.Match(regexp.MustCompile(`^[A-Za-z0-9._-]+@cbe\.com\.et$`)).Error("must be a valid @cbe.com.et email")),
		validation.Field(&r.JobTitle, validation.Required),
		validation.Field(&r.BranchCode, validation.Required),
	)
}

//before pointer
// func (r *BPSUserUpdateRequest) Validate() error {
// 	// Trim spaces
// 	r.UserID = strings.TrimSpace(r.UserID)
// 	r.FullName = strings.TrimSpace(r.FullName)
// 	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
// 	r.Email = strings.TrimSpace(r.Email)
// 	r.JobTitle = strings.TrimSpace(r.JobTitle)

// 	return validation.ValidateStruct(r,
// 		validation.Field(&r.UserID, validation.Required),
// 		validation.Field(&r.FullName, validation.Required),
// 		validation.Field(&r.PhoneNumber, validation.Required, validation.Match(phoneRegex).Error("must be a valid Ethiopian phone number (+2519xxxxxxxx or +2517xxxxxxxx)")),
// 		validation.Field(&r.Email, validation.Required, is.Email, validation.Match(regexp.MustCompile(`^[A-Za-z0-9._%+-]+@cbe\.com\.et$`)).Error("must be a valid @cbe.com.et email")),
// 		validation.Field(&r.JobTitle, validation.Required),
// 		validation.Field(&r.BranchCode, validation.Required),
// 	)
// }

type BPSUserResposenDTO struct {
	ID                bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserCode          string        `json:"user_code" bson:"user_code"` // generated
	FullName          string        `json:"full_name" bson:"full_name"`
	Username          string        `json:"username" bson:"username"`
	Email             string        `json:"email" bson:"email"`
	PhoneNumber       string        `json:"phone_number" bson:"phone_number"`
	BranchCode        []string      `json:"branch_code" bson:"branch_code"` // enum: IFB, CB
	BranchName        string        `json:"branch_name" bson:"branch_name"`
	HomeBranch        string        `json:"home_branch" bson:"home_branch"`
	Role              string        `json:"role" bson:"role"`   // enum: Maker, Checker, Aduditer
	Realm             string        `json:"realm" bson:"realm"` // default: bank
	LoginAttemptCount uint8         `json:"login_attempt_count" bson:"login_attempt_count"`
	FirstPasswordSet  bool          `json:"first_password_set" bson:"first_password_set"`
	Enabled           bool          `json:"enabled" bson:"enabled"`
	IsDeleted         bool          `json:"is_deleted" bson:"is_deleted"`
	OTPVerifyCount    uint8         `json:"otp_verfy_count" bson:"otp_verify_count"`
	OTPLastTriedAt    time.Time     `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	JobTitle          string        `json:"job_title" bson:"job_title"`
	OTPLastVerifiedAt time.Time     `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	Password          Password      `json:"login_password" bson:"login_password"`
	IsFirstTimeLogin  bool          `json:"is_first_time_login" bson:"is_first_time_login"`
	LastLoginAttempt  time.Time     `json:"last_login_attempt" bson:"last_login_attempt"`
	NextLoginAttempt  time.Time     `json:"next_login_attempt" bson:"next_login_attempt"`
	LastLogin         time.Time     `json:"last_login" bson:"last_login"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at,omitempty"`
	LastModifiedAt    time.Time     `json:"last_modifed_at" bson:"last_modifed_at,omitempty"`
	UserRole          string        `json:"user_role" bson:"user_role"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}
