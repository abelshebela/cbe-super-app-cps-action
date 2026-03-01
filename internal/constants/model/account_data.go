package model

type AccountData struct {
	AccountNumber   string `json:"AccountNumber"`
	CustomerName    string `json:"CustomerName"`
	Restriction     string `json:"Restriction"`
	Currency        string `json:"Currency"`
	WorkingBalance  string `json:"WorkingBalance"`
	CustomerID      string `json:"CustomerID"`
	AccountType     string `json:"AccountType"`
	Branch          string `json:"Branch"`
	BranchCode      string `json:"BranchCode"`
	PhoneNumber     string `json:"PhoneNo"`
	BirthOfDate     string `json:"DOB"`
	Gender          string `json:"Gender"`
	CustomerSegment string `json:"CustomerSegment"`
	RestrictionType string `json:"RestrictionType"`
	Flag            string `json:"flag,omitempty"`
	LinkedStatus    bool   `json:"linked_status"`
	ExpiryDate      string `json:"expiry_date,omitempty"`
}
