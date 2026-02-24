package core

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/model"
	"time"
)

func MapCustomerKYCToResponse(c model.CustomerKYC) dto.CustomerKYCResponse {
	return dto.CustomerKYCResponse{
		ID:           c.ID.Hex(),
		AccountType:  c.AccountType,
		CustomerCode: c.CustomerCode,
		CustomerName: dto.CustomerInfoResp{
			FirstName:   c.CustomerName.FirstName,
			MiddleName:  c.CustomerName.MiddleName,
			LastName:    c.CustomerName.LastName,
			PhoneNumber: c.CustomerName.PhoneNumber,
			Email:       c.CustomerName.Email,
			DateOfBirth: c.CustomerName.DateOfBirth,
			Gender:      c.CustomerName.Gender,
			MotherName:  c.CustomerName.MotherName,
		},
		Address: dto.AddressResp{
			Country:     c.Address.Country,
			Region:      c.Address.Region,
			City:        c.Address.City,
			SubCity:     c.Address.SubCity,
			Wereda:      c.Address.Wereda,
			Kebele:      c.Address.Kebele,
			HouseNumber: c.Address.HouseNumber,
		},
		Nationality:          c.Nationality,
		MaritalStatus:        c.MaritalStatus,
		CustomerStatus:       c.CustomerStatus,
		EmploymentStatus:     c.EmploymentStatus,
		Occupation:           c.Occupation,
		AverageMonthlyIncome: c.AverageMonthlyIncome,
		EducationStatus:      c.EducationStatus,
		SourceOfFund:         c.SourceOfFund,
		KYCStatus:            c.KYCStatus,
		LivenessCheck: dto.AlivenessCheckResp{
			IDCardFront:        c.LivenessCheck.IDCardFront,
			IDCardBack:         c.LivenessCheck.IDCardBack,
			LivenessCheckVideo: c.LivenessCheck.LivenessCheckVideo,
		},
		MoneyLaunderingFree: c.MoneyLaunderingFree,
		TermsAndConditions:  c.TermsAndConditions,
		CreatedAt:           c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           c.UpdatedAt.Format(time.RFC3339),
	}
}
