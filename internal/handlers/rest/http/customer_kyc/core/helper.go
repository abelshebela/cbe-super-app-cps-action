package core

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

func ValidateString(s string) error {
	if s == "" {
		return errors.New("string is empty")
	}

	specialCharRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	if specialCharRegex.MatchString(s) {
		return errors.New("string contains special characters")
	}
	return nil
}

func ParseRequestFromMultipleFormData(r *http.Request) (dto.CreateCustomerKYCRequest, error) {
	var req dto.CreateCustomerKYCRequest

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return req, localization.ErrorFileParseFailed
	}

	req.AccountType = constants.AccountType(r.FormValue("account_type"))
	req.CustomerName = dto.CustomerInfoRequest{
		FirstName:   r.FormValue("first_name"),
		MiddleName:  r.FormValue("middle_name"),
		LastName:    r.FormValue("last_name"),
		PhoneNumber: r.FormValue("phone_number"),
		Email:       r.FormValue("email"),
		DateOfBirth: r.FormValue("date_of_birth"),
		Gender:      r.FormValue("gender"),
		MotherName:  r.FormValue("mother_name"),
	}
	req.Address = dto.AddressRequest{
		Country:     r.FormValue("country"),
		Region:      r.FormValue("region"),
		City:        r.FormValue("city"),
		SubCity:     r.FormValue("sub_city"),
		Wereda:      r.FormValue("wereda"),
		Kebele:      r.FormValue("kebele"),
		HouseNumber: r.FormValue("house_number"),
	}
	req.Nationality = r.FormValue("nationality")
	req.MaritalStatus = constants.MaritalStatus(r.FormValue("marital_status"))
	req.EmploymentStatus = constants.EmploymentStatus(r.FormValue("employment_status"))
	req.Occupation = r.FormValue("occupation")
	req.AverageMonthlyIncome = r.FormValue("average_monthly_income")
	req.EducationStatus = r.FormValue("education_status")
	req.SourceOfFund = r.FormValue("source_of_fund")
	req.TermsAndConditions = r.FormValue("terms_and_conditions")

	faceMatchScore, _ := strconv.ParseFloat(r.FormValue("face_match_score"), 64)
	req.VerificationResult = dto.VerificationResult{
		FaceMatchScore:             faceMatchScore,
		LivenessResult:             r.FormValue("liveness_result"),
		DocumentAuthenticityResult: r.FormValue("document_authenticity_result"),
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if headers, ok := r.MultipartForm.File["id_card_front"]; ok && len(headers) > 0 {
			req.LivenessCheck.IDCardFront = headers[0]
		}
		if headers, ok := r.MultipartForm.File["id_card_back"]; ok && len(headers) > 0 {
			req.LivenessCheck.IDCardBack = headers[0]
		}
		if headers, ok := r.MultipartForm.File["liveness_video"]; ok && len(headers) > 0 {
			req.LivenessCheck.LivenessCheckVideo = headers[0]
		}
	}

	return req, nil
}
