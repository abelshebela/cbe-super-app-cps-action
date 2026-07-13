package customerkycself

type CustomerFetchDetailResponse struct {
	AccountNumber      string `xml:"AccountNo"`
	AccountName        string `xml:"AccountName"`
	Currency           string `xml:"Currency"`
	Category           string `xml:"Category"`
	AccountType        string `xml:"AccountType"`
	BranchCode         string `xml:"BranchCode"`
	BranchName         string `xml:"BranchName"`
	Balance            string `xml:"Balance"`
	RestrictionDesc    string `xml:"RestrictionDesc"`
	RestrictionType    string `xml:"RestrictionType"`
	InactiveFlag       string `xml:"InactiveFlag"`
	CustomerID         string `xml:"CustomerID"`
	CIFBranchCode      string `xml:"CIFBranchCode"`
	CIFBranchName      string `xml:"CIFBranchName"`
	CustomerName       string `xml:"CustomerName"`
	Gender             string `xml:"Gender"`
	DOB                string `xml:"DOB"`
	Mail               string `xml:"Mail"`
	Phone              string `xml:"Phone"`
	CIFRestrictType    string `xml:"CIFRestrictType"`
	CIFRestrictionDesc string `xml:"CIFRestrictionDesc"`
	CustomerGroup      string `xml:"CustomerGroup"`
	Segement           string `xml:"Segement"`
	SubSegement        string `xml:"SubSegement"`
	Industry           string `xml:"Industry"`
	Sector             string `xml:"Sector"`
	Ownership          string `xml:"Ownership"`
	Target             string `xml:"Target"`
}
