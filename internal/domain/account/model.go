package account

type AccountUser struct {
    ID              string
    PhoneNumber     string
    KYCLevel        uint8
    BranchCode      string
    FullName        string
    RegistrationType string
    AndOrStatus     bool
}


type AccountCreationResult struct {
    CustomerNumber string
    AccountNumber  string
}

type LinkedAccount struct {
    UserID            string
    CustomerNumber    string
    AccountNumber     string
    AccountHolderName string
    AccountType       string
    BranchCode        string
    RegistrationType  string
    AndOrStatus       bool
    CurrencyCode      string
    IsMain            bool
}

type MockUserData struct {
    ID            string
    PhoneNumber   string
    CustomerNumber string
}