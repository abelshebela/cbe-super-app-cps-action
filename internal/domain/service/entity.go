package service

import "time"

type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)

type ProductCodes struct {
	PRD    string
	VATPRD string
	SFPRD  string
	TRXN   string
}

type GLEntry struct {
	ProductAccount    string
	ProductBranchCode string
	ServiceAccount    string
	ServiceBranchCode string
	VatAccount        string
	VatBranchCode     string
}

type Tier struct {
	ID        string
	Min       uint64
	Max       uint64
	FeeAmount uint64
}

type Cap struct {
	KYCLevel  KYCLevel
	SingleCap uint64
	DailyCap  uint64
	MinAmount uint64
	MaxAmount uint64
}

type Service struct {
	ID                 string
	ServiceCode        string
	ServiceName        string
	ServiceType        string
	Key                string
	Cap                Cap
	CBEProductCodes    ProductCodes
	CBEIFBProductCodes ProductCodes
	AboveAmount        uint64
	AboveServiceFee    uint64
	PaymentType        string
	Tiers              []Tier
	CBEGLEntry         GLEntry
	CBEIFBGLEntry      GLEntry
	Enabled            bool
	IsDeleted          bool
	CreatedAt          time.Time
	LastModifiedAt     time.Time
	DeletedAt          time.Time
}
