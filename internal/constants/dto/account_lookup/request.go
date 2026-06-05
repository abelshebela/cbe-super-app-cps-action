package accountlookup

import (
	"cbe-super-app-cps-action/internal/constants"
)

type AccountLookUpRequest struct {
	AccountNumber string `json:"account_number,omitempty"`
	PhoneNumber   string `json:"phone_number"`
}

type CreateAccountRequest struct {
	CustomerName       string                `json:"customer_name"`
	Gender             constants.Gender      `json:"gender"`
	PhoneNumber        string                `json:"phone_number"`
	CustomerMotherName string                `json:"customer_mother_name"`
	AccountType        string                `json:"account_type"`
	AccountBranchType  constants.AccountType `json:"account_branchtype"`
	Picture            string                `json:"picture"`
	CustomerAddress    string                `json:"customer_address"`
}

type AccountCreateParams struct {
	Username           string
	Password           string
	FirstName          string
	MiddleName         string
	LastName           string
	PhoneNumber        string
	Address            string
	PostalCode         string
	ISOCountryCode     string
	AccountOffice      string
	Industry           string
	ISONationalityCode string
	ISOResidentCode    string
	UniqueID           string
	IssuesBy           string
	IssuedDate         string
	ExpiryDate         string
	Gender             string
	DateOfBirth        string
	MaritalStatus      string
	Email              string
	EmploymentStatus   string
	Occupation         string
	EmployerName       string
	EmployerAddress    string
	EmployerBusiness   string
	CustomerCurrency   string
	Salary             string
	AnnualBonus        string
	NetMonthlyIncome   string
	NetMonthlyExpence  string
	TinNumber          string
	MotherName         string
	CustomerGroup      string
}
