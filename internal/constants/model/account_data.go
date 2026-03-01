package model

type AccountData struct {
	AccountNumber   string `json:"AccountNumber"`
	CustomerName    string `json:"CustomerName"`
	Restriction     string `json:"Restriction"`
	Currency        string `json:"Currency"`
	CustomerID      string `json:"CustomerID"`
	Category        string `json:"Category"`
	AccountType     string `json:"AccountType"`
	BranchCode      string `json:"BranchCode"`
	BranchName      string `json:"BranchName"`
	DistrictName    string `json:"DistrictName"`
	PhoneNo         string `json:"PhoneNo"`
	Industry        string `json:"Industry"`
	Sector          string `json:"Sector"`
	Ownership       string `json:"Ownership"`
	CustomerSegment string `json:"CustomerSegment"`
	Target          string `json:"Target"`
	Gender          string `json:"Gender"`
	BirthOfDate     string `json:"BirthOfDate"`
	Email           string `json:"Email"`
	RestrictionType string `json:"RestrictionType"`
	PhoneNumber     string `json:"PhoneNumber"`
	Branch          string `json:"Branch"`
}
