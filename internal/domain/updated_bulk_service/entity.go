package updatedbulkservice

import "time"

type Tier struct {
	Min       float64 `json:"min" bson:"min"`
	Max       float64 `json:"max" bson:"max"`
	FeeAmount float64 `json:"feeAmount" bson:"feeAmount"`
}

type ProductCodes struct {
	PRD    string `json:"PRD" bson:"PRD"`
	VATPRD string `json:"VATPRD" bson:"VATPRD"`
	SFPRD  string `json:"SFPRD" bson:"SFPRD"`
	TRXN   string `json:"TRXN" bson:"TRXN"`
}

type CBglEntry struct {
	CBglProductAccount    string `json:"CBglProductAccount" bson:"CBglProductAccount"`
	CBglProductBranchcode string `json:"CBglProductBranchcode" bson:"CBglProductBranchcode"`
	CBglServiceAccount    string `json:"CBglServiceAccount" bson:"CBglServiceAccount"`
	CBglServiceBranchcode string `json:"CBglServiceBranchcode" bson:"CBglServiceBranchcode"`
	CBglVatAccount        string `json:"CBglVatAccount" bson:"CBglVatAccount"`
	CBglVatBranchcode     string `json:"CBglVatBranchcode" bson:"CBglVatBranchcode"`
}

type IFBglEntry struct {
	IFBglProductAccount    string `json:"IFBglProductAccount" bson:"IFBglProductAccount"`
	IFBglProductBranchcode string `json:"IFBglProductBranchcode" bson:"IFBglProductBranchcode"`
	IFBglServiceAccount    string `json:"IFBglServiceAccount" bson:"IFBglServiceAccount"`
	IFBglServiceBranchcode string `json:"IFBglServiceBranchcode" bson:"IFBglServiceBranchcode"`
	IFBglVatAccount        string `json:"IFBglVatAccount" bson:"IFBglVatAccount"`
	IFBglVatBranchcode     string `json:"IFBglVatBranchcode" bson:"IFBglVatBranchcode"`
}

type ServiceDetails struct {
	Key               string       `json:"key" bson:"key"`
	ServiceCode       string       `json:"serviceCode" bson:"serviceCode"`
	ServiceName       string       `json:"serviceName" bson:"serviceName"`
	SingleCap         float64      `json:"singleCap" bson:"singleCap"`
	DailyCap          float64      `json:"dailyCap" bson:"dailyCap"`
	MinAmount         float64      `json:"minAmount" bson:"minAmount"`
	ServiceType       string       `json:"serviceType" bson:"serviceType"` // "bill" | "transfer"
	PaymentType       string       `json:"paymentType" bson:"paymentType"` // "tier percentage" | "flat amount"
	VAT               float64      `json:"vat" bson:"vat"`
	Tiers             []Tier       `json:"tiers" bson:"tiers"`
	MinAmountVIRTUAL  float64      `json:"minAmountVIRTUAL" bson:"minAmountVIRTUAL"`
	SingleCapLevelOne float64      `json:"singleCapLevelOne" bson:"singleCapLevelOne"`
	DailyCapLevelOne  float64      `json:"dailyCapLevelOne" bson:"dailyCapLevelOne"`
	AboveAmount       float64      `json:"aboveAmount" bson:"aboveAmount"`
	AboveServiceFee   float64      `json:"aboveServiceFee" bson:"aboveServiceFee"`
	Prefix            string       `json:"prefix,omitempty" bson:"prefix,omitempty"`
	USSDEnabled       bool         `json:"USSDEnabled" bson:"USSDEnabled"`
	Enabled           bool         `json:"enabled" bson:"enabled"`
	IsDeleted         bool         `json:"isDeleted" bson:"isDeleted"`
	DateCreated       time.Time    `json:"dateCreated" bson:"dateCreated"`
	LastModified      time.Time    `json:"lastModified" bson:"lastModified"`
	CreatedDate       time.Time    `json:"createdDate" bson:"createdDate"`
	ProductCodes      ProductCodes `json:"productCodes" bson:"productCodes"`
	IFBProductCodes   ProductCodes `json:"IFBproductCodes" bson:"IFBproductCodes"`
	CBglEntry         CBglEntry    `json:"CBglEntry" bson:"CBglEntry"`
	IFBglEntry        IFBglEntry   `json:"IFBglEntry" bson:"IFBglEntry"`
}
